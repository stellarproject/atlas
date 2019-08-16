package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/ehazlett/atlas/ds/filters"
)

func (s *Server) Lookup(ctx context.Context, req *api.LookupRequest) (*api.LookupResponse, error) {
	// TODO: enable filters by req.Query
	nf := &filters.Name{Value: req.Query}
	records, err := s.ds.Search(nf)
	if err != nil {
		return nil, err
	}

	return &api.LookupResponse{
		Records: records,
	}, nil
}
