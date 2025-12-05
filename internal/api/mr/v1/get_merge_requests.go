package mr_v1

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	api "github.com/vlanse/glmr/internal/pb/mr/v1"
	"github.com/vlanse/glmr/internal/service/mr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetMergeRequests(ctx context.Context, req *api.GetMergeRequestsRequest) (*api.GetMergeRequestsResponse, error) {
	mrg, err := s.mrSvc.GetMergeRequests(ctx, mr.Filter{
		SkipApprovedByMe: req.GetFilter().GetSkipApprovedByMe(),
		ButStillShowMine: req.GetFilter().GetButStillShowMine(),
		ShowOnlyMine:     req.GetFilter().GetShowOnlyMine(),
		DoNotShowDrafts:  req.GetFilter().GetDoNotShowDrafts(),
	})
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	res := &api.GetMergeRequestsResponse{
		Groups: lo.Map(mrg, func(item mr.MergeRequestsGroup, _ int) *api.GetMergeRequestsResponse_Group {
			return &api.GetMergeRequestsResponse_Group{
				Name: item.GroupName,
				Summary: &api.GetMergeRequestsResponse_Group_Summary{
					Total:   int32(item.Summary.Total),
					Overdue: int32(item.Summary.Overdue),
				},
				MergeRequests: lo.Map(
					item.MergeRequests,
					func(item mr.MergeRequest, _ int) *api.GetMergeRequestsResponse_MergeRequest {
						return &api.GetMergeRequestsResponse_MergeRequest{
							Description: item.Description,
							Iid:         item.IID,
							Project: &api.GetMergeRequestsResponse_MergeRequest_Project{
								Id:   item.Project.ID,
								Name: item.Project.Name,
								Url:  item.Project.WebURL,
							},
							Url: item.URL,
							Author: &api.GetMergeRequestsResponse_MergeRequest_User{
								Username:  item.Author.Username,
								AvatarUrl: item.Author.AvatarURL,
								Url:       item.Author.WebURL,
								IsMe:      item.Author.IsMe,
							},
							ApprovedBy: lo.Map(item.Approvals, func(item mr.Approval, _ int) *api.GetMergeRequestsResponse_MergeRequest_User {
								return &api.GetMergeRequestsResponse_MergeRequest_User{
									Username:  item.User.Username,
									AvatarUrl: item.User.AvatarURL,
									Trusted:   item.User.IsOwner,
									Url:       item.User.WebURL,
									IsMe:      item.User.IsMe,
								}
							}),
							Status: &api.GetMergeRequestsResponse_MergeRequest_Status{
								PipelineFailed:  item.Status.PipelineFailed,
								Conflict:        item.Status.Conflict,
								Ready:           item.Status.Ready,
								Outdated:        item.Status.Outdated,
								Pending:         item.Status.Pending,
								EditorAvailable: s.editorSvc.IsProjectConfigured(item.Project.ID),
							},
							Comments: &api.GetMergeRequestsResponse_MergeRequest_Comments{
								ResolvedCount:   int32(item.CommentStats.ResolvedCount),
								UnresolvedCount: int32(item.CommentStats.UnresolvedCount),
							},
							Age:            fmt.Sprintf("%dd", int(time.Since(item.CreatedAt).Hours()/24)),
							ApprovedBefore: item.ApprovedBefore,
							Issues: lo.Map(item.Issues, func(item mr.Issue, _ int) *api.GetMergeRequestsResponse_MergeRequest_Issue {
								return &api.GetMergeRequestsResponse_MergeRequest_Issue{
									Key: item.Key,
									Url: item.URL,
								}
							}),
							DiffStatsSummary: &api.GetMergeRequestsResponse_MergeRequest_DiffStatsSummary{
								Additions: item.DiffStatsSummary.Additions,
								Deletions: item.DiffStatsSummary.Deletions,
								FileCount: item.DiffStatsSummary.FileCount,
							},
						}
					},
				),
			}
		}),
	}

	return res, nil
}
