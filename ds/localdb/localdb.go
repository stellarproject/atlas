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
	"os"
	"path/filepath"

	"github.com/ehazlett/atlas/ds"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketID = ds.ServiceID + ".v1"
)

// LocalDB is a BoltDB backed datastore
type LocalDB struct {
	db *bolt.DB
}

// New returns a new datastore
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

// ID returns the id of the service
func (l *LocalDB) ID() string {
	return "localdb"
}

// Close is used to close and release all resources
func (l *LocalDB) Close() error {
	if l.db != nil {
		l.db.Close()
	}
	return nil
}
