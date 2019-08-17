package localdb

import (
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/ds"
	bolt "go.etcd.io/bbolt"
)

func (l *LocalDB) Search(query string, filters ...ds.Filter) ([]*api.Record, error) {
	var records []*api.Record
	if err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketID))
		if err := b.ForEach(func(k, v []byte) error {
			if !matchKey(query, string(k)) {
				return nil
			}
			var r []*api.Record
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			records = append(records, r...)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return records, nil
	}

	// apply filters
	var res []*api.Record
	for _, f := range filters {
		r, err := f.Apply(records)
		if err != nil {
			return nil, err
		}
		res = append(res, r...)
	}
	return res, nil
}

func matchKey(query, domain string) bool {
	if query == "*" || query == domain {
		return true
	}

	// host.foo.invalid, *.foo.invalid
	if strings.Index(domain, "*") == 0 {
		s := getRootDomain(query)
		d := getRootDomain(domain)
		return s == d
	}
	return false
}

func getRootDomain(v string) string {
	x := strings.Split(v, ".")
	if len(x) < 3 {
		return v
	}
	return fmt.Sprintf("%s.%s", x[len(x)-2], x[len(x)-1])
}
