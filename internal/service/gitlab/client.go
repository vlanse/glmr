package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/vlanse/glmr/internal/util/request"
)

const (
	tokenHeader = "Private-Token"
	authHeader  = "Authorization"
)

type client struct {
	mx       sync.RWMutex
	settings Settings
}

func newClient(settings Settings) *client {
	return &client{
		settings: settings,
	}
}

func (c *client) updateSettings(settings Settings) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.settings = settings
}

func (c *client) getSettings() Settings {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.settings
}

func (c *client) getProject(ctx context.Context, projectID int64) (Project, error) {
	settings := c.getSettings()

	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d", settings.BaseURL, projectID)),
		map[string]string{
			tokenHeader: settings.Token,
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

func (c *client) getStarredProjects(ctx context.Context, userID int64) ([]Project, error) {
	settings := c.getSettings()

	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/users/%d/starred_projects", settings.BaseURL, userID)),
		map[string]string{
			tokenHeader: settings.Token,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get starred projects from gitlab: %w", err)
	}

	var res []Project
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal starred projects info: %w", err)
	}

	return res, nil
}

func (c *client) getApprovalRules(ctx context.Context, projectID int64) ([]ApprovalRule, error) {
	settings := c.getSettings()

	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/projects/%d/approval_rules", settings.BaseURL, projectID)),
		map[string]string{
			tokenHeader: settings.Token,
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
	settings := c.getSettings()

	data, err := request.GET(
		ctx,
		request.MustURL(fmt.Sprintf("%s/api/v4/user", settings.BaseURL)),
		map[string]string{
			tokenHeader: settings.Token,
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

func (c *client) getMergeRequestDiscussions(ctx context.Context, projectID, mergeRequestIID int64) ([]Discussion, error) {
	settings := c.getSettings()

	res := make([]Discussion, 0)
	const perPage = 100
	for page := 1; ; page++ {
		data, err := request.GET(
			ctx,
			request.MustURL(
				fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/discussions", settings.BaseURL, projectID, mergeRequestIID),
				"per_page", strconv.Itoa(perPage), "page", strconv.Itoa(page),
			),
			map[string]string{
				tokenHeader: settings.Token,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get merge request discussions from gitlab: %w", err)
		}

		pageRes := make([]Discussion, 0)
		if err = json.Unmarshal(data, &pageRes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal merge request discussions info: %w", err)
		}

		res = append(res, pageRes...)
		if len(pageRes) < perPage {
			break
		}
	}

	return res, nil
}

func (c *client) getProjectMergeRequestsGQ(ctx context.Context, projectFullPath string) ([]MergeRequestGQ, error) {
	q := struct {
		Query string `json:"query"`
	}{
		Query: fmt.Sprintf(projectMRsRequest, projectFullPath),
	}
	reqData, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal projects merge request: %w", err)
	}

	settings := c.getSettings()

	data, err := request.POST(ctx,
		request.MustURL(fmt.Sprintf("%s/api/graphql", settings.BaseURL)),
		map[string]string{
			authHeader:     fmt.Sprintf("Bearer %s", settings.Token),
			"Content-Type": "application/json",
		},
		reqData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get project merge requests from gitlab: %w", err)
	}

	res := struct {
		Data struct {
			Project struct {
				MergeRequests struct {
					Nodes []MergeRequestGQ `json:"nodes"`
				} `json:"mergeRequests"`
			} `json:"project"`
		} `json:"data"`
	}{}
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal merge requests info: %w", err)
	}

	return res.Data.Project.MergeRequests.Nodes, nil
}
