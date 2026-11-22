package gitlab

import "context"

type Settings struct {
	BaseURL string
	Token   string
}

type Service struct {
	cl *client
}

func NewService(settings Settings) *Service {
	return &Service{
		cl: newClient(settings),
	}
}

func (s *Service) UpdateSettings(settings Settings) {
	s.cl.updateSettings(settings)
}

func (s *Service) GetBaseURL() string {
	return s.cl.getSettings().BaseURL
}

func (s *Service) GetProject(ctx context.Context, projectID int64) (Project, error) {
	return s.cl.getProject(ctx, projectID)
}

func (s *Service) GetStarredProjects(ctx context.Context, userID int64) ([]Project, error) {
	return s.cl.getStarredProjects(ctx, userID)
}

func (s *Service) GetApprovalRules(ctx context.Context, projectID int64) ([]ApprovalRule, error) {
	return s.cl.getApprovalRules(ctx, projectID)
}

func (s *Service) GetCurrentUser(ctx context.Context) (User, error) {
	return s.cl.getCurrentUser(ctx)
}

func (s *Service) GetMergeRequestDiscussions(ctx context.Context, projectID, mergeRequestIID int64) ([]Discussion, error) {
	return s.cl.getMergeRequestDiscussions(ctx, projectID, mergeRequestIID)
}

func (s *Service) GetProjectMergeRequestsGQ(ctx context.Context, projectPath string) ([]MergeRequestGQ, error) {
	return s.cl.getProjectMergeRequestsGQ(ctx, projectPath)
}
