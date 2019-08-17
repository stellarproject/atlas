package localdb

import "testing"

func TestGetRootDomain(t *testing.T) {
	expected := "foo.invalid"
	d := getRootDomain("foo.invalid")

	if d != expected {
		t.Fatalf("expected %s; received %s", expected, d)
	}
}

func TestGetRootDomainSubdomain(t *testing.T) {
	expected := "foo.invalid"
	d := getRootDomain("host.foo.invalid")

	if d != expected {
		t.Fatalf("expected %s; received %s", expected, d)
	}
}

func TestGetRootDomainMultipleSubdomain(t *testing.T) {
	expected := "foo.invalid"
	d := getRootDomain("foo.bar.host.foo.invalid")

	if d != expected {
		t.Fatalf("expected %s; received %s", expected, d)
	}
}

func TestMatchKeyExact(t *testing.T) {
	key := "foo.invalid"
	domain := "foo.invalid"

	if !matchKey(key, domain) {
		t.Fatalf("expected key %s to match %s", key, domain)
	}
}

func TestMatchKeySubdomain(t *testing.T) {
	key := "host.foo.invalid"
	domain := "*.foo.invalid"

	if !matchKey(key, domain) {
		t.Fatalf("expected key %s to match %s", key, domain)
	}
}

func TestMatchKeySubdomainInvalid(t *testing.T) {
	key := "host.foo.invalid.local"
	domain := "*.foo.invalid"

	if matchKey(key, domain) {
		t.Fatalf("expected miss %s to match %s", key, domain)
	}
}

func TestMatchKeyWildcard(t *testing.T) {
	key := "*"
	domain := "foo.invalid"

	if !matchKey(key, domain) {
		t.Fatalf("expected key %s to match %s", key, domain)
	}
}
