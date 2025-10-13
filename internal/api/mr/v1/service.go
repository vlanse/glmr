package mr_v1

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/vlanse/glmr/internal/pb/mr/v1"
	"github.com/vlanse/glmr/internal/service/mr"
	"google.golang.org/grpc"
)

type Service struct {
	api.UnsafeMergeRequestsServer

	mrSvc *mr.Service
}

func New(mrSvc *mr.Service) *Service {
	return &Service{
		mrSvc: mrSvc,
	}
}

func (s *Service) Register(
	ctx context.Context,
	srv *grpc.Server,
	mux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption,
) error {
	api.RegisterMergeRequestsServer(srv, s)
	return api.RegisterMergeRequestsHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
