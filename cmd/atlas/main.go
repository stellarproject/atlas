/*
   Copyright 2019 Stellar Project

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
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/stellarproject/atlas/version"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.BuildVersion()
	app.Author = "@stellarproject"
	app.Email = ""
	app.Usage = version.Description
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug logging",
		},
		cli.StringFlag{
			Name:   "port, p",
			Usage:  "port for the DNS service",
			Value:  "53",
			EnvVar: "ATLAS_PORT",
		},
		cli.StringFlag{
			Name:   "redis-url, r",
			Usage:  "uri for datastore backend",
			Value:  "redis://127.0.0.1:6379/0",
			EnvVar: "ATLAS_REDIS_URL",
		},
		cli.StringFlag{
			Name:   "address, a",
			Usage:  "grpc address",
			Value:  "tcp://127.0.0.1:9000",
			EnvVar: "ATLAS_GRPC_ADDR",
		},
		cli.StringFlag{
			Name:   "upstream-dns",
			Usage:  "upstream dns server",
			Value:  "9.9.9.9",
			EnvVar: "ATLAS_UPSTREAM_DNS",
		},
		cli.StringFlag{
			Name:   "dnsmasq-conf, c",
			Usage:  "path for generated dnsmasq config",
			Value:  "/etc/dnsmasq.conf",
			EnvVar: "ATLAS_DNSMASQ_CONF",
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
