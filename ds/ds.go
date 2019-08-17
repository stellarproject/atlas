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

package ds

import (
	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

const (
	// ServiceID is the id of the datastore service
	ServiceID = "com.evanhazlett.atlas.datastore"
)

// Filter allows for filtering of records
type Filter interface {
	// Apply is the implementation needed to filter records
	Apply(r []*api.Record) ([]*api.Record, error)
}

// Datastore defines the datastore interface
type Datastore interface {
	// ID returns the id of the datastore
	ID() string
	// Get gets the specified records by key
	Get(key string) ([]*api.Record, error)
	// Set sets the key to the records
	Set(key string, v []*api.Record) error
	// Search returns a list of records optionally filtered
	Search(key string, filters ...Filter) ([]*api.Record, error)
	// Delete deletes records by key
	Delete(key string) error
	// Close optionally closes any resources in use by the datastore
	Close() error
}
