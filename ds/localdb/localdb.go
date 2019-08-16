package localdb

import (
	"os"
	"path/filepath"

	"github.com/ehazlett/atlas/ds"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketID = ds.ServiceID + ".v1"
)

type LocalDB struct {
	db *bolt.DB
}

func New(dbPath string) (*LocalDB, error) {
	baseDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "error creating database path %s", baseDir)
	}
	db, err := bolt.Open(dbPath, 0664, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening database %s", dbPath)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketID)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &LocalDB{
		db: db,
	}, nil
}

func (l *LocalDB) ID() string {
	return "localdb"
}

func (l *LocalDB) Close() error {
	if l.db != nil {
		l.db.Close()
	}
	return nil
}
