package localdb

import (
	"encoding/json"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	bolt "go.etcd.io/bbolt"
)

func (l *LocalDB) Set(key string, r []*api.Record) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return l.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketID))
		err := b.Put([]byte(key), data)
		return err
	})
}
