package atlas

import (
	"net/url"

	"github.com/ehazlett/atlas/ds"
	"github.com/ehazlett/atlas/ds/localdb"
	"github.com/pkg/errors"
)

var (
	// ErrUnsupportedDatastore is the error returned when an unsupported datastore is specified
	ErrUnsupportedDatastore = errors.New("unsupported datastore")
)

// GetDatastore returns an instance of a specific datastore specified by the URI
func GetDatastore(dbURI string) (ds.Datastore, error) {
	// TODO: parse db uri and return proper datastore
	u, err := url.Parse(dbURI)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "localdb":
		return localdb.New(u.Path)
	default:
		return nil, errors.Wrap(err, u.Scheme)
	}
}
