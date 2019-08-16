package localdb

import (
	bolt "go.etcd.io/bbolt"
)

func (l *LocalDB) Delete(key string) error {
	return l.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketID))
		err := b.Delete([]byte(key))
		return err
	})
}
