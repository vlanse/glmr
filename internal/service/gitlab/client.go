package gitlab

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vlanse/glmr/internal/util/request"
)

const (
	tokenHeader = "Private-Token"
)

type client struct {
	baseURL string
	token   string
}

func newClient(baseURL string, token string) *client {
	return &client{
		baseURL: baseURL,
		token:   token,
	}
}

func (c *client) getProjectMergeRequests(ctx context.Context, projectID int64) ([]MergeRequest, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(
			fmt.Sprintf("%s/api/v4/projects/%d/merge_requests", c.baseURL, projectID),
			"per_page", "100", "page", "1", "state", "opened",
		),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge requests from gitlab: %w", err)
	}

	res := make([]MergeRequest, 0)
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merge requests info: %w", err)
	}

	return res, nil
}

func (c *client) getProject(ctx context.Context, projectID int64) (Project, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d", c.baseURL, projectID)),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return Project{}, fmt.Errorf("failed to get project info from gitlab: %w", err)
	}

	var res Project
	if err = json.Unmarshal(data, &res); err != nil {
		return Project{}, fmt.Errorf("failed to unmarshal project info: %w", err)
	}

	return res, nil

}

func (c *client) getApprovalRules(ctx context.Context, projectID int64) ([]ApprovalRule, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d/approval_rules", c.baseURL, projectID)),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get project approval rules from gitlab: %w", err)
	}

	res := make([]ApprovalRule, 0)
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project approval rules: %w", err)
	}

	return res, nil
}

func (c *client) getCurrentUser(ctx context.Context) (User, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/user", c.baseURL)),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return User{}, fmt.Errorf("failed to get current user info from gitlab: %w", err)
	}

	var res User
	if err = json.Unmarshal(data, &res); err != nil {
		return User{}, fmt.Errorf("failed to unmarshal current user info: %w", err)
	}

	return res, nil
}

func (c *client) getMergeRequestsApprovals(ctx context.Context, projectID, mrIID int64) (Approval, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/approvals", c.baseURL, projectID, mrIID)),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return Approval{}, fmt.Errorf("failed to get MR approvals info from gitlab: %w", err)
	}

	var res Approval
	if err = json.Unmarshal(data, &res); err != nil {
		return Approval{}, fmt.Errorf("failed to unmarshal MR approval info rules: %w", err)
	}

	return res, nil
}

func (c *client) getMergeRequestDiscussions(ctx context.Context, projectID, mergeRequestIID int64) ([]Discussion, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(
			fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/discussions", c.baseURL, projectID, mergeRequestIID),
			"per_page", "100", "page", "1", "state", "opened",
		),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request discussions from gitlab: %w", err)
	}

	res := make([]Discussion, 0)
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merge request discussions info: %w", err)
	}

	return res, nil
}

func (c *client) getMergeRequestCommits(ctx context.Context, projectID, mergeRequestIID int64) ([]Commit, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(
			fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/commits", c.baseURL, projectID, mergeRequestIID),
			"per_page", "100", "page", "1",
		),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get merge request commits from gitlab: %w", err)
	}

	res := make([]Commit, 0)
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merge request commits info: %w", err)
	}

	return res, nil
}

func (c *client) getMergeRequestInfo(ctx context.Context, projectID, mergeRequestIID int64) (MergeRequestInfo, error) {
	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d", c.baseURL, projectID, mergeRequestIID)),
		map[string]string{
			tokenHeader: c.token,
		},
	)
	if err != nil {
		return MergeRequestInfo{}, fmt.Errorf("failed to get merge request info from gitlab: %w", err)
	}

	var res MergeRequestInfo
	if err = json.Unmarshal(data, &res); err != nil {
		return MergeRequestInfo{}, fmt.Errorf("failed to unmarshal merge request info: %w", err)
	}

	return res, nil
}
