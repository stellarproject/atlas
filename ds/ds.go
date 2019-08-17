package ds

import (
	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

const (
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