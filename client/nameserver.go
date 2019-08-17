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
package client

import (
	"context"
	"fmt"
	"strings"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func (c *Client) Create(name string, records []*api.Record) error {
	ctx := context.Background()
	if _, err := c.nameserverService.Create(ctx, &api.CreateRequest{
		Name:    name,
		Records: records,
	}); err != nil {
		return err
	}
	return nil
}

func (c *Client) Lookup(query string) ([]*api.Record, error) {
	ctx := context.Background()
	resp, err := c.nameserverService.Lookup(ctx, &api.LookupRequest{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) List() ([]*api.Record, error) {
	ctx := context.Background()
	resp, err := c.nameserverService.List(ctx, &api.ListRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Records, nil
}

func (c *Client) Delete(name string) error {
	ctx := context.Background()
	if _, err := c.nameserverService.Delete(ctx, &api.DeleteRequest{
		Name: name,
	}); err != nil {
		return err
	}
	return nil
}

func (c *Client) RecordType(rtype string) (api.RecordType, error) {
	switch strings.ToUpper(rtype) {
	case "A":
		return api.RecordType_A, nil
	case "CNAME":
		return api.RecordType_CNAME, nil
	case "SRV":
		return api.RecordType_SRV, nil
	case "TXT":
		return api.RecordType_TXT, nil
	case "MX":
		return api.RecordType_MX, nil
	default:
		return api.RecordType_UNKNOWN, fmt.Errorf("unsupported record type %q", rtype)
	}
}
