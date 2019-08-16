package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/ehazlett/atlas"
	"github.com/ehazlett/atlas/server"
	"github.com/ehazlett/atlas/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func runServer(cx *cli.Context) error {
	cfg := &atlas.Config{
		BindAddress:     cx.String("bind"),
		Datastore:       cx.String("datastore"),
		GRPCAddress:     cx.String("address"),
		UpstreamDNSAddr: cx.String("upstream-dns"),
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"version": version.Version,
		"commit":  version.GitCommit,
	}).Infof("starting %s", version.Name)
	if err := srv.Start(); err != nil {
		return err
	}

	// create grpc server
	grpcOpts, err := getGRPCOptions(cfg)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(grpcOpts...)

	// register atlas server
	if err := srv.Register(grpcServer); err != nil {
		return err
	}

	proto, ep, err := getGRPCEndpoint(cfg.GRPCAddress)
	if err != nil {
		return err
	}
	l, err := net.Listen(proto, ep)
	if err != nil {
		return err
	}
	defer l.Close()

	logrus.WithField("addr", cfg.GRPCAddress).Debug("starting grpc server")
	go grpcServer.Serve(l)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	doneCh := make(chan bool, 1)
	go func() {
		for {
			select {
			case sig := <-signals:
				switch sig {
				case syscall.SIGUSR1:
					logrus.Debug("generating debug profile")
					profilePath, err := srv.GenerateProfile()
					if err != nil {
						logrus.Error(err)
						continue
					}
					logrus.WithFields(logrus.Fields{
						"profile": profilePath,
					}).Info("generated memory profile")
				case syscall.SIGTERM, syscall.SIGINT:
					logrus.Info("shutting down")
					if err := srv.Stop(); err != nil {
						logrus.Error(err)
					}
					doneCh <- true
				default:
					logrus.Warnf("unhandled signal %s", sig)
				}
			}
		}
	}()

	<-doneCh

	return nil
}

func getGRPCEndpoint(addr string) (string, string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", "", err
	}
	// only tcp/unix are allowed
	var ep string
	switch u.Scheme {
	case "tcp":
		ep = u.Host
	case "unix":
		ep = u.Path
	default:
		return "", "", fmt.Errorf("unsupported grpc listener protocol: %s", u.Scheme)
	}

	return u.Scheme, ep, nil
}

func getGRPCOptions(cfg *atlas.Config) ([]grpc.ServerOption, error) {
	grpcOpts := []grpc.ServerOption{}
	if cfg.TLSServerCertificate != "" && cfg.TLSServerKey != "" {
		logrus.WithFields(logrus.Fields{
			"cert": cfg.TLSServerCertificate,
			"key":  cfg.TLSServerKey,
		}).Debug("configuring TLS for GRPC")
		cert, err := tls.LoadX509KeyPair(cfg.TLSServerCertificate, cfg.TLSServerKey)
		if err != nil {
			return nil, err

		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates:       []tls.Certificate{cert},
			ClientAuth:         tls.RequestClientCert,
			InsecureSkipVerify: cfg.TLSInsecureSkipVerify,
		})
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}
	return grpcOpts, nil
}
