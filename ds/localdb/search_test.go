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
	expected := "bar.host.foo.invalid"
	d := getRootDomain("foo.bar.host.foo.invalid")

	if d != expected {
		t.Fatalf("expected %s; received %s", expected, d)
	}
}

func TestMatchKeyExact(t *testing.T) {
	query := "foo.invalid"
	domain := "foo.invalid"

	if !matchKey(query, domain) {
		t.Fatalf("expected query %s to match %s", query, domain)
	}
}

func TestMatchKeySubdomain(t *testing.T) {
	query := "host.foo.invalid"
	domain := "*.foo.invalid"

	if !matchKey(query, domain) {
		t.Fatalf("expected query %s to match %s", query, domain)
	}
}

func TestMatchKeySubdomainInvalid(t *testing.T) {
	query := "host.foo.invalid.local"
	domain := "*.foo.invalid"

	if matchKey(query, domain) {
		t.Fatalf("expected miss %s to match %s", query, domain)
	}
}

func TestMatchKeySubdomainMulti(t *testing.T) {
	query := "host.bar.foo.invalid"
	domain := "*.bar.foo.invalid"

	if !matchKey(query, domain) {
		t.Fatalf("expected query %s to match %s", query, domain)
	}
}

func TestMatchKeySubdomainMultiRoot(t *testing.T) {
	query := "host.bar.foo.invalid"
	domain := "*.inv.foo.invalid"

	if matchKey(query, domain) {
		t.Fatalf("expected query %s to not match %s", query, domain)
	}
}

func TestMatchKeyWildcard(t *testing.T) {
	query := "*"
	domain := "foo.invalid"

	if !matchKey(query, domain) {
		t.Fatalf("expected query %s to match %s", query, domain)
	}
}
