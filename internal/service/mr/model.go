package mr

import "time"

type Filter struct {
	SkipApprovedByMe bool
	ButStillShowMine bool
	ShowOnlyMine     bool
	DoNotShowDrafts  bool
}

type User struct {
	ID        int64 `json:"id"`
	Username  string
	AvatarURL string
	IsOwner   bool
	WebURL    string
	IsMe      bool
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
	WebURL string
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

type Issue struct {
	Key string
	URL string
}

type DiffStatsSummary struct {
	Additions int64
	Deletions int64
	FileCount int64
}

type PluginResult struct {
	HTML      string
	PlainText string
}

type MergeRequest struct {
	IID              int64 // "short" gitlab ID
	Project          Project
	CreatedAt        time.Time
	Description      string
	URL              string
	Draft            bool
	Author           User
	Approvals        []Approval
	Commits          []Commit
	Pipeline         Pipeline
	Discussions      []Discussion
	CommentStats     CommentStats
	Status           Status
	ApprovedBefore   bool
	Issues           []Issue
	DiffStatsSummary DiffStatsSummary
}

type ApprovalRule struct {
	Name  string
	Users []User
}

type Project struct {
	ID                int64
	Name              string
	GroupName         string
	WebURL            string
	PathWithNamespace string
	MergeRequests     []MergeRequest
	ApprovalRules     []ApprovalRule
	Plugins           []PluginResult
}

type Summary struct {
	Total          int
	Visible        int
	Overdue        int
	OverdueVisible int
	Draft          int
	DraftVisible   int
}
type MergeRequestsGroup struct {
	GroupName     string
	MergeRequests []MergeRequest
	Summary       Summary
}
