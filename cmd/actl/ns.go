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
	Name:  "delete",
	Usage: "delete nameserver record",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Usage: "resource record type (A, CNAME, TXT, SRV, MX)",
			Value: "A",
		},
	},
	ArgsUsage: "<NAME>",
	Action: func(c *cli.Context) error {
		client, err := getClient(c)
		if err != nil {
			return err
		}
		defer client.Close()

		t := c.String("type")
		name := c.Args().First()

		if name == "" {
			return fmt.Errorf("you must enter a name")
		}

		if err := client.Delete(t, name); err != nil {
			return err
		}

		fmt.Printf("removed %s (%s)\n", name, t)

		return nil
	},
}
