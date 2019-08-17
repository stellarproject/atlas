package filters

import (
	"testing"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func TestFilterRecordType(t *testing.T) {
	records := []*api.Record{
		{
			Type: api.RecordType_A,
			Name: "foo.invalid",
		},
		{
			Type: api.RecordType_CNAME,
			Name: "foo.invalid",
		},
	}

	f := &RecordType{
		Type: api.RecordType_A,
	}

	res, err := f.Apply(records)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 result; received %d", len(res))
	}
}
