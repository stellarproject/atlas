package localdb

import (
	"encoding/json"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	bolt "go.etcd.io/bbolt"
)

func (l *LocalDB) Get(key string) ([]*api.Record, error) {
	var records []*api.Record
	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketID))
		v := b.Get([]byte(key))

		return json.Unmarshal(v, records)
	})

	return records, err
}
