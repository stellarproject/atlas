/*
   Copyright 2019 Evan Hazlett <ejhazlett@gmail.com>

   Permission is hereby granted, free of charge, to any person obtaining a copy of
   this software and associated documentation files (the "Software"), to deal in the
   Software without restriction, including without limitation the rights to use, copy,
   modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
   and to permit persons to whom the Software is furnished to do so, subject to the
   following conditions:

   The above copyright notice and this permission notice shall be included in all copies
   or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
   PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
   FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
   TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE
   USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

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
		MetricsAddr:     cx.String("metrics-addr"),
		CacheTTL:        cx.Duration("cache-ttl"),
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

	waitForExit(srv)

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
