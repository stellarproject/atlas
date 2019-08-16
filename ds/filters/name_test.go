package filters

import (
	"testing"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func TestNameFilter(t *testing.T) {
	name := "foo.invalid"

	f := &Name{
		Value: name,
	}

	r := []*api.Record{
		{
			Type: api.RecordType_A,
			Name: name,
		},
	}

	if !f.Match(r) {
		t.Fatalf("expected to match %s", name)
	}
}

func TestNameFilterSubdomain(t *testing.T) {
	query := "host.foo.invalid"
	name := "*.foo.invalid"

	f := &Name{
		Value: query,
	}

	r := []*api.Record{
		{
			Type: api.RecordType_A,
			Name: name,
		},
	}

	if !f.Match(r) {
		t.Fatalf("expected %s to match subdomain %s", query, name)
	}
}

func TestNameFilterSubdomainWithRoot(t *testing.T) {
	query := "foo.invalid"
	name := "*.foo.invalid"

	f := &Name{
		Value: query,
	}

	r := []*api.Record{
		{
			Type: api.RecordType_A,
			Name: name,
		},
	}

	if !f.Match(r) {
		t.Fatalf("expected %s to match subdomain %s", query, name)
	}
}
