package gitlab

import "context"

type Service struct {
	cl *client
}

func NewService(baseURL string, token string) *Service {
	return &Service{
		cl: newClient(baseURL, token),
	}
}

func (s *Service) GetProject(ctx context.Context, projectID int64) (Project, error) {
	return s.cl.getProject(ctx, projectID)
}

func (s *Service) GetProjectMergeRequests(ctx context.Context, projectID int64) ([]MergeRequest, error) {
	return s.cl.getProjectMergeRequests(ctx, projectID)
}

func (s *Service) GetApprovalRules(ctx context.Context, projectID int64) ([]ApprovalRule, error) {
	return s.cl.getApprovalRules(ctx, projectID)
}

func (s *Service) GetCurrentUser(ctx context.Context) (User, error) {
	return s.cl.getCurrentUser(ctx)
}

func (s *Service) GetMergeRequestApprovals(ctx context.Context, projectID int64, mergeRequestIID int64) (Approval, error) {
	return s.cl.getMergeRequestsApprovals(ctx, projectID, mergeRequestIID)
}

func (s *Service) GetMergeRequestDiscussions(ctx context.Context, projectID, mergeRequestIID int64) ([]Discussion, error) {
	return s.cl.getMergeRequestDiscussions(ctx, projectID, mergeRequestIID)
}

func (s *Service) GetMergeRequestCommits(ctx context.Context, projectID, mergeRequestIID int64) ([]Commit, error) {
	return s.cl.getMergeRequestCommits(ctx, projectID, mergeRequestIID)
}

func (s *Service) GetMergeRequestInfo(ctx context.Context, projectID, mergeRequestIID int64) (MergeRequestInfo, error) {
	return s.cl.getMergeRequestInfo(ctx, projectID, mergeRequestIID)
}
