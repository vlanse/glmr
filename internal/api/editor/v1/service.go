package editor_v1

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/vlanse/glmr/internal/pb/editor/v1"
	"github.com/vlanse/glmr/internal/service/editor"
	"google.golang.org/grpc"
)

type Service struct {
	api.UnsafeEditorServer

	editorSvc *editor.Service
}

func New(editorSvc *editor.Service) *Service {
	return &Service{
		editorSvc: editorSvc,
	}
}

func (s *Service) Register(
	ctx context.Context,
	srv *grpc.Server,
	mux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption,
) error {
	api.RegisterEditorServer(srv, s)
	return api.RegisterEditorHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
