package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	ptypes "github.com/gogo/protobuf/types"
)

func (s *Server) Delete(ctx context.Context, req *api.DeleteRequest) (*ptypes.Empty, error) {
	if err := s.ds.Delete(req.Name); err != nil {
		return empty, err
	}
	return empty, nil
}
