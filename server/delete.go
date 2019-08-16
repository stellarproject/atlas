package server

import (
	"context"

	api "github.com/ehazlett/atlas/api/services/nameserver/v1"
	ptypes "github.com/gogo/protobuf/types"
)

func (s *Server) Delete(ctx context.Context, req *api.DeleteRequest) (*ptypes.Empty, error) {
	return empty, nil
}
