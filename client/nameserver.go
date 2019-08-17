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
