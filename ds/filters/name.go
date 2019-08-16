package filters

import (
	"strings"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

// Name filter matches exactly on a name or a subdomain
type Name struct {
	Value string
}

// Match matches on an exact name or subdomain including wildcards
func (n *Name) Match(records []*api.Record) bool {
	for _, r := range records {
		if r.Name == n.Value {
			return true
		}
		// check for subdomain
		if n.matchSubdomain(r.Name) {
			return true
		}
	}
	return false
}

func (n *Name) matchSubdomain(v string) bool {
	x := strings.Split(n.Value, ".")
	var domain string
	if len(x) == 2 {
		domain = strings.Join(x[0:len(x)], ".")
	} else {
		domain = strings.Join(x[1:2], ".")
	}
	parts := strings.SplitN(v, domain, 2)

	if len(parts) != 2 {
		return false
	}

	host := strings.Trim(parts[0], ".")

	// check for wildcard
	if host == "*" {
		return true
	}

	return host == v
}
