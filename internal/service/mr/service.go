package mr

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond/v2"
	"github.com/samber/lo"
	"github.com/vlanse/glmr/internal/service/gitlab"
)

const (
	poolWorkerCount = 30

	ownerRuleName = "Owner"

	pipelineFailedStatus  = "failed"
	pipelineSuccessStatus = "success"
)

var jiraIssueRx = regexp.MustCompile(`\[([A-Z]+-\d+)]`)

type Service struct {
	settings  Settings
	gitlabSvc *gitlab.Service
	pool      pond.Pool

	currentUser *User
	dataMx      sync.Mutex

	projectsByID map[int64]gitlab.Project
}

func NewService(gitlabSvc *gitlab.Service) *Service {
	return &Service{
		gitlabSvc:    gitlabSvc,
		pool:         pond.NewPool(poolWorkerCount),
		projectsByID: make(map[int64]gitlab.Project),
	}
}

func (s *Service) UpdateSettings(settings Settings) {
	s.dataMx.Lock()
	defer s.dataMx.Unlock()
	s.settings = settings
}

func (s *Service) GetMergeRequests(ctx context.Context, filter Filter) ([]MergeRequestsGroup, error) {
	var currentUserName string
	if err := func() error {
		s.dataMx.Lock()
		defer s.dataMx.Unlock()
		if s.currentUser == nil {
			user, err := s.gitlabSvc.GetCurrentUser(ctx)
			if err != nil {
				return fmt.Errorf("could not get current user information: %w", err)
			}
			s.currentUser = &User{
				Username:  user.Username,
				AvatarURL: user.AvatarURL,
				WebURL:    user.WebURL,
			}
		}
		currentUserName = s.currentUser.Username
		return nil
	}(); err != nil {
		return nil, err
	}

	projects := s.settings.GetProjects()

	var err error

	if projects, err = s.enrichProjectInfoGQ(ctx, projects); err != nil {
		return nil, err
	}

	if projects, err = s.enrichProjectMRDiscussions(ctx, projects); err != nil {
		return nil, err
	}

	// if projects, err = s.enrichProjectInfo(ctx, projects); err != nil {
	// 	return nil, err
	// }
	//
	// if projects, err = s.enrichProjectMRInfo(ctx, projects); err != nil {
	// 	return nil, err
	// }

	projects = s.fillIssues(projects)

	projects = fillOwners(projects)

	projects = sortApprovers(currentUserName, projects)

	projects = fillStatuses(projects)

	projects = countDiscussions(projects)

	projects = setApprovedBefore(currentUserName, projects)

	groups := groupMergeRequests(projects)

	groups = fillGroupSummaries(groups)

	groups = filterMergeRequests(groups, currentUserName, filter)

	for _, g := range groups {
		sort.SliceStable(g.MergeRequests, func(i, j int) bool {
			return g.MergeRequests[i].CreatedAt.Before(g.MergeRequests[j].CreatedAt)
		})
	}

	return groups, nil
}

func (s *Service) fillIssues(projects []Project) []Project {
	if len(s.settings.JIRA.URL) == 0 {
		return projects
	}
	for i, project := range projects {
		for j, mr := range project.MergeRequests {
			issueKeys := jiraIssueRx.FindAllStringSubmatch(mr.Description, -1)
			var keys []string
			for _, issueKey := range issueKeys {
				if len(issueKey) < 2 {
					continue
				}
				key := issueKey[1]
				if lo.Contains(keys, key) {
					continue
				}
				keys = append(keys, key)
				projects[i].MergeRequests[j].Issues = append(projects[i].MergeRequests[j].Issues, Issue{
					Key: key,
					URL: fmt.Sprintf("%s/browse/%s", s.settings.JIRA.URL, key),
				})
			}
		}
	}
	return projects
}

func fillGroupSummaries(groups []MergeRequestsGroup) []MergeRequestsGroup {
	for i, g := range groups {
		for _, mr := range g.MergeRequests {
			groups[i].Summary.Total++
			if mr.Status.Outdated {
				groups[i].Summary.Overdue++
			}
		}
	}
	return groups
}

func filterMergeRequests(groups []MergeRequestsGroup, currentUsername string, filter Filter) []MergeRequestsGroup {
	for i, group := range groups {
		groups[i].MergeRequests = lo.Filter(group.MergeRequests, func(item MergeRequest, _ int) bool {
			stillShowMine := filter.ButStillShowMine && item.Author.Username == currentUsername

			if filter.DoNotShowDrafts && strings.HasPrefix(item.Description, "Draft:") {
				if !stillShowMine {
					return false
				}
			}

			if filter.SkipApprovedByMe {
				if lo.ContainsBy(item.Approvals, func(item Approval) bool {
					return item.User.Username == currentUsername
				}) {
					if !stillShowMine {
						return false
					}
				}
			}

			if filter.ShowOnlyMine && item.Author.Username != currentUsername {
				return false
			}
			return true
		})
	}
	return groups
}

func fillOwners(projects []Project) []Project {
	ownersByProjectID := make(map[int64][]string)
	for _, project := range projects {
		if owners, found := lo.Find(project.ApprovalRules, func(item ApprovalRule) bool {
			return item.Name == ownerRuleName
		}); found {
			ownersByProjectID[project.ID] = append(ownersByProjectID[project.ID],
				lo.Map(owners.Users, func(item User, _ int) string {
					return item.Username
				})...,
			)
		}

		for i, mr := range project.MergeRequests {
			for j, a := range mr.Approvals {
				if idx := lo.IndexOf(ownersByProjectID[project.ID], a.User.Username); idx != -1 {
					project.MergeRequests[i].Approvals[j].User.IsOwner = true
				}
			}
		}
	}
	return projects
}

func fillStatuses(projects []Project) []Project {
	for i, project := range projects {
		for j, mr := range project.MergeRequests {
			status := Status{
				PipelineFailed: mr.Pipeline.Status == pipelineFailedStatus,
				Outdated:       time.Since(mr.CreatedAt) > time.Hour*24*10,
				Conflict:       mr.Status.Conflict,
			}
			status.Ready = mr.Pipeline.Status == pipelineSuccessStatus && !status.Conflict

			status.Pending = !lo.Contains(
				[]string{pipelineSuccessStatus, pipelineFailedStatus}, mr.Pipeline.Status,
			) && !status.Conflict

			projects[i].MergeRequests[j].Status = status
		}
	}
	return projects
}

func countDiscussions(projects []Project) []Project {
	for i, project := range projects {
		for j, mr := range project.MergeRequests {
			var resolved, unresolved int
			for _, d := range mr.Discussions {
				for _, n := range d.Notes {
					if !n.Resolvable {
						continue
					}
					if n.Resolved {
						resolved++
					} else {
						unresolved++
					}
				}
			}
			projects[i].MergeRequests[j].CommentStats = CommentStats{
				ResolvedCount:   resolved,
				UnresolvedCount: unresolved,
			}
		}
	}
	return projects
}

func setApprovedBefore(currentUserName string, projects []Project) []Project {
	for i, project := range projects {
		for j, mr := range project.MergeRequests {
			if mr.Author.Username == currentUserName {
				continue
			}
			if _, alreadyApproved := lo.Find(mr.Approvals, func(item Approval) bool {
				return item.User.Username == currentUserName
			}); alreadyApproved {
				continue
			}
			for _, d := range mr.Discussions {
				for _, n := range d.Notes {
					if n.Author.Username != currentUserName {
						continue
					}
					if strings.Contains(n.Body, "approved this merge request") {
						projects[i].MergeRequests[j].ApprovedBefore = true
					}
				}
			}
		}
	}
	return projects
}

func sortApprovers(currentUserName string, projects []Project) []Project {
	for i, project := range projects {
		for j, mr := range project.MergeRequests {
			sort.SliceStable(mr.Approvals, func(i, j int) bool {
				a1, a2 := mr.Approvals[i], mr.Approvals[j]
				if a1.User.Username == currentUserName && a2.User.Username != currentUserName {
					return true
				}
				if a1.User.Username != currentUserName && a2.User.Username == currentUserName {
					return false
				}

				if a1.User.IsOwner && !a2.User.IsOwner {
					return true
				}
				if !a1.User.IsOwner && a2.User.IsOwner {
					return false
				}

				return a1.User.Username < a2.User.Username
			})
			projects[i].MergeRequests[j].Approvals = mr.Approvals
		}
	}
	return projects
}

func groupMergeRequests(projects []Project) []MergeRequestsGroup {
	groupProjects := lo.PartitionBy(projects, func(item Project) string {
		return item.GroupName
	})

	var res []MergeRequestsGroup
	for _, ps := range groupProjects {
		mrg := MergeRequestsGroup{
			GroupName: ps[0].GroupName,
		}
		for _, project := range ps {
			mrg.MergeRequests = append(mrg.MergeRequests, project.MergeRequests...)
		}
		res = append(res, mrg)
	}
	return res
}
