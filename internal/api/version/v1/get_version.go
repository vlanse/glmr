package version_v1

import (
	"context"

	api "github.com/vlanse/glmr/internal/pb/version/v1"
	"github.com/vlanse/glmr/internal/util/version"
)

func (s *Service) GetVersion(ctx context.Context, _ *api.GetVersionRequest) (*api.GetVersionResponse, error) {
	updateVersion, updateMessage, err := version.CheckForUpdates(ctx)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	res := &api.GetVersionResponse{
		CurrentVersion: version.GetCurrent(),
		Update: &api.GetVersionResponse_Update{
			Version:      updateVersion,
			ReleaseNotes: updateMessage,
			Error:        errMsg,
		},
	}
	return res, nil
}
