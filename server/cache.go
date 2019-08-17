package server

import (
	"encoding/json"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func (s *Server) cacheRecords(name string, records []*api.Record) error {
	if s.cfg.CacheTTL == 0 {
		return nil
	}

	data, err := json.Marshal(records)
	if err != nil {
		return err
	}
	if err := s.cache.Set(name, string(data)); err != nil {
		return err
	}
	return nil
}

func (s *Server) fromCache(name string) ([]*api.Record, error) {
	if s.cfg.CacheTTL == 0 {
		return nil, nil
	}

	kv := s.cache.Get(name)
	if kv == nil {
		return nil, nil
	}

	var records []*api.Record
	d := []byte(kv.Value.(string))
	if err := json.Unmarshal(d, &records); err != nil {
		return nil, err
	}
	return records, nil
}
