package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	"github.com/sirupsen/logrus"
)

func (s *Server) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	// TODO: enable filters by req.Query
	logrus.Debug("searching records")
	records, err := s.ds.Search()
	if err != nil {
		return &api.ListResponse{}, err
	}

	return &api.ListResponse{
		Records: records,
	}, nil
}
