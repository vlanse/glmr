package version_v1

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/vlanse/glmr/internal/pb/version/v1"
	"google.golang.org/grpc"
)

type Service struct {
	api.UnsafeVersionServer
}

func New() *Service {
	return &Service{}
}

func (s *Service) Register(
	ctx context.Context,
	srv *grpc.Server,
	mux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption,
) error {
	api.RegisterVersionServer(srv, s)
	return api.RegisterVersionHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
