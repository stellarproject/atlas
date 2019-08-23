/*
   Copyright 2019 Stellar Project

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

package server

import (
	"encoding/json"

	api "github.com/stellarproject/atlas/api/services/nameserver/v1"
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
