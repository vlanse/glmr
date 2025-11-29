package editor_v1

import (
	"context"

	api "github.com/vlanse/glmr/internal/pb/editor/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) OpenProject(ctx context.Context, req *api.OpenProjectRequest) (*api.OpenProjectResponse, error) {
	if err := s.editorSvc.OpenProject(ctx, req.GetProjectId()); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	return nil, nil
}
