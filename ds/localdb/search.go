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

import (
	"encoding/json"
	"strings"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/ds"
	bolt "go.etcd.io/bbolt"
)

// Search returns records for the specified query using optional filters
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
	return strings.Join(x[1:], ".")
}
