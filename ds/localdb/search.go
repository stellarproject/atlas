package localdb

import (
	"encoding/json"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/ds"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func (l *LocalDB) Search(filters ...ds.Filter) ([]*api.Record, error) {
	var records []*api.Record
	if err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketID))
		if err := b.ForEach(func(k, v []byte) error {
			logrus.Debugf("k=%s v=%s", string(k), string(v))
			var r []*api.Record
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			for _, f := range filters {
				if !f.Match(r) {
					return nil
				}
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
	return records, nil
}
