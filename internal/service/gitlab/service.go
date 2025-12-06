package gitlab

import "context"

type Service struct {
	cl *client
}

func NewService(baseURL, token string) *Service {
	return &Service{
		cl: newClient(baseURL, token),
	}
}

func (s *Service) UpdateSettings(baseURL, token string) {
	s.cl.baseURL = baseURL
	s.cl.token = token
}

func (s *Service) GetBaseURL() string {
	return s.cl.baseURL
}

func (s *Service) GetProject(ctx context.Context, projectID int64) (Project, error) {
	return s.cl.getProject(ctx, projectID)
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
