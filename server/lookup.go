package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
)

func (s *Server) Lookup(ctx context.Context, req *api.LookupRequest) (*api.LookupResponse, error) {
	// TODO: enable filters
	records, err := s.ds.Search(req.Query)
	if err != nil {
		return nil, err
	}

	return &api.LookupResponse{
		Records: records,
	}, nil
}
