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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/containerd/typeurl"
	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/api/types"
	"github.com/urfave/cli"
)

var listRecordsCommand = cli.Command{
	Name:  "list",
	Usage: "list nameserver records",
	Action: func(c *cli.Context) error {
		client, err := getClient(c)
		if err != nil {
			return err
		}
		defer client.Close()

		records, err := client.List()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		fmt.Fprintf(w, "NAME\tTYPE\tVALUE\tOPTIONS\n")
		for _, r := range records {
			opts := ""
			if r.Options != nil {
				v, err := typeurl.UnmarshalAny(r.Options)
				if err != nil {
					return err
				}
				if o, ok := v.(types.NameserverOption); ok {
					opts = o.String()
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Name, r.Type, r.Value, opts)
		}
		w.Flush()

		return nil
	},
}

var createRecordCommand = cli.Command{
	Name:  "create",
	Usage: "create nameserver record",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Usage: "resource record type (A, CNAME, TXT, SRV, MX)",
			Value: "A",
		},
		// TODO: handle resource record options
	},
	ArgsUsage: "<NAME> <VALUE>",
	Action: func(c *cli.Context) error {
		client, err := getClient(c)
		if err != nil {
			return err
		}
		defer client.Close()

		t := c.String("type")
		name := c.Args().First()
		value := c.Args().Get(1)

		if name == "" || value == "" {
			return fmt.Errorf("you must enter a name and value")
		}

		rType, err := client.RecordType(t)
		if err != nil {
			return err
		}
		record := &api.Record{
			Type:  rType,
			Name:  name,
			Value: value,
		}

		if err := client.Create(name, []*api.Record{record}); err != nil {
			return err
		}

		fmt.Printf("added %s=%s (%s)\n", name, value, t)

		return nil
	},
}

var deleteRecordCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete nameserver record",
	Flags:     []cli.Flag{},
	ArgsUsage: "<NAME>",
	Action: func(c *cli.Context) error {
		client, err := getClient(c)
		if err != nil {
			return err
		}
		defer client.Close()

		name := c.Args().First()

		if name == "" {
			return fmt.Errorf("you must enter a name")
		}

		if err := client.Delete(name); err != nil {
			return err
		}

		fmt.Printf("removed %s\n", name)

		return nil
	},
}
