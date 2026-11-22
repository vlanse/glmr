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

func (c *client) getProject(ctx context.Context, projectID int64) (Project, error) {
	settings := c.getSettings()

	return fetchJSON[Project](ctx, "project info",
		request.WithQuery(fmt.Sprintf("%s/api/v4/projects/%d", settings.BaseURL, projectID)),
		map[string]string{
			tokenHeader: settings.Token,
		},
	)
}

func (c *client) getStarredProjects(ctx context.Context, userID int64) ([]Project, error) {
	settings := c.getSettings()

	return fetchJSON[[]Project](ctx, "starred projects",
		request.WithQuery(fmt.Sprintf("%s/api/v4/users/%d/starred_projects", settings.BaseURL, userID)),
		map[string]string{
			tokenHeader: settings.Token,
		},
	)
}

func (c *client) getApprovalRules(ctx context.Context, projectID int64) ([]ApprovalRule, error) {
	settings := c.getSettings()

	return fetchJSON[[]ApprovalRule](ctx, "project approval rules",
		request.WithQuery(fmt.Sprintf("%s/api/v4/projects/%d/approval_rules", settings.BaseURL, projectID)),
		map[string]string{
			tokenHeader: settings.Token,
		},
	)
}

func (c *client) getCurrentUser(ctx context.Context) (User, error) {
	settings := c.getSettings()

	return fetchJSON[User](ctx, "current user info",
		request.WithQuery(fmt.Sprintf("%s/api/v4/user", settings.BaseURL)),
		map[string]string{
			tokenHeader: settings.Token,
		},
	)
}

func (c *client) getMergeRequestDiscussions(ctx context.Context, projectID, mergeRequestIID int64) ([]Discussion, error) {
	settings := c.getSettings()

	res := make([]Discussion, 0)
	const perPage = 100
	for page := 1; ; page++ {
		pageRes, err := fetchJSON[[]Discussion](ctx, "merge request discussions",
			request.WithQuery(
				fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d/discussions", settings.BaseURL, projectID, mergeRequestIID),
				"per_page", strconv.Itoa(perPage), "page", strconv.Itoa(page),
			),
			map[string]string{
				tokenHeader: settings.Token,
			},
		)
		if err != nil {
			return nil, err
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
		request.WithQuery(fmt.Sprintf("%s/api/graphql", settings.BaseURL)),
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

func (c *client) getSettings() Settings {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.settings
}

func fetchJSON[T any](ctx context.Context, what, url string, headers map[string]string) (T, error) {
	var res T

	data, err := request.GET(ctx, url, headers)
	if err != nil {
		return res, fmt.Errorf("failed to get %s from gitlab: %w", what, err)
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return res, fmt.Errorf("failed to unmarshal %s: %w", what, err)
	}

	return res, nil
}
