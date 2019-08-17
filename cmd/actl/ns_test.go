package main

import (
	"fmt"
	"testing"
)

func TestParseRecord(t *testing.T) {
	tt := "a"
	v := "127.0.0.1"
	q := fmt.Sprintf("%s:%s", tt, v)

	r, err := parseRecord(q)
	if err != nil {
		t.Fatal(err)
	}

	if r.Type != tt {
		t.Fatalf("expected type %s; received %s", tt, r.Type)
	}
	if r.Value != v {
		t.Fatalf("expected value %s; received %s", v, r.Value)
	}
}
