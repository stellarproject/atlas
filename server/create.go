package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	ptypes "github.com/gogo/protobuf/types"
)

func (s *Server) Create(ctx context.Context, req *api.CreateRequest) (*ptypes.Empty, error) {
	if err := s.ds.Set(req.Name, req.Records); err != nil {
		return empty, err
	}
	return empty, nil
}
