package filters

import (
	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

// RecordType is a filter that limits the records to the defined RecordType
type RecordType struct {
	// Type is the RecordType on which to filter
	Type api.RecordType
}

// Apply filters records by the specified type
func (f *RecordType) Apply(records []*api.Record) ([]*api.Record, error) {
	var res []*api.Record
	for _, r := range records {
		if r.Type != f.Type {
			continue
		}

		res = append(res, r)
	}
	return res, nil
}
