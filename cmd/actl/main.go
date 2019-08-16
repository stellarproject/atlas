package main

import (
	"fmt"
	"os"

	"github.com/ehazlett/atlas"
	"github.com/ehazlett/atlas/client"
	"github.com/ehazlett/atlas/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()
	app.Name = "actl"
	app.Version = version.BuildVersion()
	app.Author = "@ehazlett"
	app.Email = ""
	app.Usage = version.Description
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug logging",
		},
		cli.StringFlag{
			Name:  "addr, a",
			Usage: "atlas grpc address",
			Value: "tcp://127.0.0.1:9000",
		},
		cli.StringFlag{
			Name:  "cert, c",
			Usage: "atlas client certificate",
			Value: "",
		},
		cli.StringFlag{
			Name:  "key, k",
			Usage: "atlas client key",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "skip-verify",
			Usage: "skip TLS verification",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}
	app.Commands = []cli.Command{
		listRecordsCommand,
		createRecordCommand,
		deleteRecordCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func getClient(c *cli.Context) (*client.Client, error) {
	cert := c.GlobalString("cert")
	key := c.GlobalString("key")
	skipVerification := c.GlobalBool("skip-verify")

	cfg := &atlas.Config{
		TLSClientCertificate:  cert,
		TLSClientKey:          key,
		TLSInsecureSkipVerify: skipVerification,
	}

	opts, err := client.DialOptionsFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	opts = append(opts,
		grpc.WithBlock(),
		grpc.WithUserAgent(fmt.Sprintf("%s/%s", version.Name, version.Version)),
	)

	addr := c.GlobalString("addr")
	return client.NewClient(addr, opts...)
}
