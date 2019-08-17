package main

import (
	"os"

	"github.com/ehazlett/atlas/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.BuildVersion()
	app.Author = "@ehazlett"
	app.Email = ""
	app.Usage = version.Description
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug logging",
		},
		cli.StringFlag{
			Name:  "bind, b",
			Usage: "bind address for the DNS service",
			Value: "udp://0.0.0.0:53",
		},
		cli.StringFlag{
			Name:  "datastore, d",
			Usage: "uri for datastore backend",
			Value: "localdb:///etc/atlas/atlas.db",
		},
		cli.StringFlag{
			Name:  "address, a",
			Usage: "grpc address",
			Value: "tcp://127.0.0.1:9000",
		},
		cli.StringFlag{
			Name:  "upstream-dns",
			Usage: "upstream dns server",
			Value: "9.9.9.9:53",
		},
		cli.DurationFlag{
			Name:  "cache-ttl",
			Usage: "builtin cache ttl (default: disabled)",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		return nil
	}
	app.Action = runServer

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
