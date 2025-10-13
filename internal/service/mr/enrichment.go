package mr

import (
	"context"
	"fmt"
	"sync"

	"github.com/samber/lo"
	"github.com/vlanse/glmr/internal/service/gitlab"
)

func (s *Service) enrichProjectInfo(ctx context.Context, projects []Project) ([]Project, error) {
	allProjectIDs := lo.Map(projects, func(item Project, _ int) int64 {
		return item.ID
	})

	var mx sync.Mutex
	mrByProject := make(map[int64][]gitlab.MergeRequest, len(projects))
	rulesByProject := make(map[int64][]gitlab.ApprovalRule, len(projects))

	group := s.pool.NewGroup()
	for _, projectID := range allProjectIDs {
		group.SubmitErr(
			func() error {
				mr, err := s.gitlabSvc.GetProjectMergeRequests(ctx, projectID)
				if err != nil {
					return err
				}
				mx.Lock()
				defer mx.Unlock()
				mrByProject[projectID] = mr
				return nil
			},
			func() error {
				rules, err := s.gitlabSvc.GetApprovalRules(ctx, projectID)
				if err != nil {
					return err
				}
				mx.Lock()
				defer mx.Unlock()
				rulesByProject[projectID] = rules
				return nil
			},
		)
	}
	if err := group.Wait(); err != nil {
		return nil, fmt.Errorf("collect projects: %w", err)
	}

	for i, p := range projects {
		projects[i].MergeRequests = lo.Map(mrByProject[p.ID], func(mr gitlab.MergeRequest, _ int) MergeRequest {
			return MergeRequest{
				IID:         mr.IID,
				Project:     p,
				Description: mr.Title,
				Author: User{
					Username:  mr.Author.Username,
					AvatarURL: mr.Author.AvatarURL,
				},
				CreatedAt: mr.CreatedAt,
				URL:       mr.WebURL,
				Status: Status{
					Conflict: mr.HasConflicts,
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
					}
				}),
			}
		})
	}

	return projects, nil
}

func (s *Service) enrichProjectMRInfo(ctx context.Context, projects []Project) ([]Project, error) {
	var mx sync.Mutex

	approvalsByMR := make(map[int64]map[int64][]gitlab.ApprovedBy, len(projects))
	mrInfoByMR := make(map[int64]map[int64]gitlab.MergeRequestInfo, len(projects))
	discussionsByMR := make(map[int64]map[int64][]gitlab.Discussion, len(projects))
	commitsByMR := make(map[int64]map[int64][]gitlab.Commit, len(projects))

	group := s.pool.NewGroup()
	for _, project := range projects {
		for _, mr := range project.MergeRequests {

			group.SubmitErr(
				func() error {
					approvals, err := s.gitlabSvc.GetMergeRequestApprovals(ctx, project.ID, mr.IID)
					if err != nil {
						return err
					}
					mx.Lock()
					defer mx.Unlock()
					approvalsByMR[project.ID] = lo.Assign(
						approvalsByMR[project.ID], map[int64][]gitlab.ApprovedBy{mr.IID: approvals.ApprovedBy},
					)
					return nil
				},
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
				func() error {
					commits, err := s.gitlabSvc.GetMergeRequestCommits(ctx, project.ID, mr.IID)
					if err != nil {
						return err
					}
					mx.Lock()
					defer mx.Unlock()
					commitsByMR[project.ID] = lo.Assign(
						commitsByMR[project.ID], map[int64][]gitlab.Commit{mr.IID: commits},
					)
					return nil
				},
				func() error {
					info, err := s.gitlabSvc.GetMergeRequestInfo(ctx, project.ID, mr.IID)
					if err != nil {
						return err
					}
					mx.Lock()
					defer mx.Unlock()
					mrInfoByMR[project.ID] = lo.Assign(
						mrInfoByMR[project.ID], map[int64]gitlab.MergeRequestInfo{mr.IID: info},
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
			projects[i].MergeRequests[j].Approvals = lo.Map(
				approvalsByMR[p.ID][mr.IID], func(a gitlab.ApprovedBy, _ int) Approval {
					return Approval{
						User: User{
							Username:  a.User.Username,
							AvatarURL: a.User.AvatarURL,
						},
						ApprovedAt: a.ApprovedAt,
					}
				},
			)
			projects[i].MergeRequests[j].Commits = lo.Map(
				commitsByMR[p.ID][mr.IID], func(c gitlab.Commit, _ int) Commit {
					return Commit{
						AuthorName:  c.AuthorName,
						AuthorEmail: c.AuthorEmail,
						CreatedAt:   c.CreatedAt,
					}
				},
			)

			projects[i].MergeRequests[j].Pipeline.Status = mrInfoByMR[p.ID][mr.IID].Pipeline.Status

			projects[i].MergeRequests[j].Discussions = lo.Map(
				discussionsByMR[p.ID][mr.IID], func(d gitlab.Discussion, _ int) Discussion {
					return Discussion{
						Notes: lo.Map(d.Notes, func(item gitlab.Note, _ int) Note {
							return Note{
								Author: User{
									Username:  item.Author.Username,
									AvatarURL: item.Author.AvatarURL,
								},
								ResolvedBy: User{
									Username:  item.ResolvedBy.Username,
									AvatarURL: item.ResolvedBy.AvatarURL,
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
