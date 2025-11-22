package mr

import "time"

type Filter struct {
	SkipApprovedByMe bool
	ButStillShowMine bool
	ShowOnlyMine     bool
	DoNotShowDrafts  bool
}

type User struct {
	Username  string
	AvatarURL string
	IsOwner   bool
	WebURL    string
}

type Approval struct {
	User       User
	ApprovedAt time.Time
}

type Commit struct {
	AuthorName  string
	AuthorEmail string
	CreatedAt   time.Time
}

type Note struct {
	Author     User
	ResolvedBy User
	Resolvable bool
	Resolved   bool
	CreatedAt  time.Time
	ResolveAt  time.Time
	Body       string
}

type Discussion struct {
	Notes []Note
}

type Pipeline struct {
	Status string
}

type CommentStats struct {
	ResolvedCount   int
	UnresolvedCount int
}

type Status struct {
	Conflict       bool
	PipelineFailed bool
	Ready          bool
	Outdated       bool
	Pending        bool
}

type MergeRequest struct {
	IID            int64 // "short" gitlab ID
	Project        Project
	CreatedAt      time.Time
	Description    string
	URL            string
	Author         User
	Approvals      []Approval
	Commits        []Commit
	Pipeline       Pipeline
	Discussions    []Discussion
	CommentStats   CommentStats
	Status         Status
	ApprovedBefore bool
}

type ApprovalRule struct {
	Name  string
	Users []User
}

type Project struct {
	ID            int64
	Name          string
	GroupName     string
	WebURL        string
	MergeRequests []MergeRequest
	ApprovalRules []ApprovalRule
}

type Summary struct {
	Total   int
	Overdue int
}
type MergeRequestsGroup struct {
	GroupName     string
	MergeRequests []MergeRequest
	Summary       Summary
}
