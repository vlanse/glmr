package mr

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/vlanse/glmr/internal/service/gitlab"
)

func (s *Service) enrichProjectMRDiscussions(ctx context.Context, projects []Project) ([]Project, error) {
	var mx sync.Mutex
	discussionsByMR := make(map[int64]map[int64][]gitlab.Discussion, len(projects))

	group := s.pool.NewGroup()
	for _, project := range projects {
		for _, mr := range project.MergeRequests {
			group.SubmitErr(
				func() error {
					discussions, err := s.gitlabSvc.GetMergeRequestDiscussions(ctx, project.ID, mr.IID)
					if err != nil {
						return err
					}
					mx.Lock()
					defer mx.Unlock()
					discussionsByMR[project.ID] = lo.Assign(
						discussionsByMR[project.ID], map[int64][]gitlab.Discussion{mr.IID: discussions},
					)
					return nil
				},
			)
		}
	}

	if err := group.Wait(); err != nil {
		return nil, fmt.Errorf("enrich merge requests from gitlab: %w", err)
	}

	for i, p := range projects {
		for j, mr := range projects[i].MergeRequests {
			projects[i].MergeRequests[j].Discussions = lo.Map(
				discussionsByMR[p.ID][mr.IID], func(d gitlab.Discussion, _ int) Discussion {
					return Discussion{
						Notes: lo.Map(d.Notes, func(item gitlab.Note, _ int) Note {
							return Note{
								Author: User{
									Username:  item.Author.Username,
									AvatarURL: item.Author.AvatarURL,
									WebURL:    item.Author.WebURL,
									IsMe:      s.currentUser.Username == item.Author.Username,
								},
								ResolvedBy: User{
									Username:  item.ResolvedBy.Username,
									AvatarURL: item.ResolvedBy.AvatarURL,
									WebURL:    item.ResolvedBy.WebURL,
									IsMe:      s.currentUser.Username == item.ResolvedBy.Username,
								},
								Resolved:   item.Resolved,
								Resolvable: item.Resolvable,
								CreatedAt:  item.CreatedAt,
								ResolveAt:  item.ResolvedAt,
								Body:       item.Body,
							}
						}),
					}
				},
			)
		}
	}

	return projects, nil
}

func (s *Service) enrichProjectInfoGQ(ctx context.Context, projects []Project) ([]Project, error) {
	allProjectIDs := lo.Map(projects, func(item Project, _ int) int64 {
		return item.ID
	})

	var mx sync.Mutex
	group := s.pool.NewGroup()
	for _, projectID := range allProjectIDs {
		group.SubmitErr(
			func() error {
				if found := func() bool {
					mx.Lock()
					defer mx.Unlock()
					_, found := s.projectsByID[projectID]
					return found
				}(); found {
					return nil
				}

				project, err := s.gitlabSvc.GetProject(ctx, projectID)
				if err != nil {
					return err
				}
				mx.Lock()
				defer mx.Unlock()
				s.projectsByID[projectID] = project
				return nil
			},
		)
	}
	if err := group.Wait(); err != nil {
		return nil, fmt.Errorf("collect projects: %w", err)
	}

	mrByProject := make(map[int64][]gitlab.MergeRequestGQ, len(projects))
	rulesByProject := make(map[int64][]gitlab.ApprovalRule, len(projects))
	group = s.pool.NewGroup()
	for _, p := range s.projectsByID {
		group.SubmitErr(
			func() error {
				mrs, err := s.gitlabSvc.GetProjectMergeRequestsGQ(ctx, p.PathWithNamespace)
				if err != nil {
					return err
				}
				mx.Lock()
				defer mx.Unlock()
				mrByProject[p.ID] = mrs
				return nil

			},
			func() error {
				rules, err := s.gitlabSvc.GetApprovalRules(ctx, p.ID)
				if err != nil {
					return err
				}
				mx.Lock()
				defer mx.Unlock()
				rulesByProject[p.ID] = rules
				return nil
			},
		)
	}

	if err := group.Wait(); err != nil {
		return nil, fmt.Errorf("collect projects MRs: %w", err)
	}

	for i, p := range projects {
		p.WebURL = s.projectsByID[p.ID].WebURL

		projects[i].MergeRequests = lo.Map(mrByProject[p.ID], func(mr gitlab.MergeRequestGQ, _ int) MergeRequest {
			return MergeRequest{
				IID:         mr.IID,
				Project:     p,
				CreatedAt:   mr.CreatedAt,
				Description: mr.Title,
				URL:         mr.WebURL,
				Author: User{
					Username:  mr.Author.Username,
					AvatarURL: s.fixURL(mr.Author.AvatarURL),
					WebURL:    mr.Author.WebURL,
					IsMe:      s.currentUser.Username == mr.Author.Username,
				},
				Approvals: lo.Map(mr.ApprovedBy.Nodes, func(item gitlab.UserGQ, _ int) Approval {
					return Approval{
						User: User{
							Username:  item.Username,
							AvatarURL: s.fixURL(item.AvatarURL),
							WebURL:    item.WebURL,
							IsMe:      s.currentUser.Username == item.Username,
						},
					}
				}),
				Pipeline: Pipeline{
					Status: strings.ToLower(mr.HeadPipeline.Status),
				},
				Status: Status{
					Conflict: mr.Conflicts,
				},
				DiffStatsSummary: DiffStatsSummary{
					Additions: mr.DiffStatsSummary.Additions,
					Deletions: mr.DiffStatsSummary.Deletions,
					FileCount: mr.DiffStatsSummary.FileCount,
				},
			}
		})

		projects[i].ApprovalRules = lo.Map(rulesByProject[p.ID], func(r gitlab.ApprovalRule, _ int) ApprovalRule {
			return ApprovalRule{
				Name: r.Name,
				Users: lo.Map(r.EligibleApprovers, func(item gitlab.User, _ int) User {
					return User{
						Username:  item.Username,
						AvatarURL: item.AvatarURL,
						WebURL:    item.WebURL,
						IsMe:      s.currentUser.Username == item.Username,
					}
				}),
			}
		})
	}

	return projects, nil
}

func (s *Service) fixURL(url string) string {
	if strings.HasPrefix(url, "/") {
		return s.gitlabSvc.GetBaseURL() + url
	}
	return url
}
