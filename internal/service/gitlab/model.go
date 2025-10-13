package gitlab

import "time"

type MergeRequest struct {
	ID           int64     `json:"id"`
	IID          int64     `json:"iid"`
	ProjectID    int64     `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	WebURL       string    `json:"web_url"`
	HasConflicts bool      `json:"has_conflicts"`
	Author       struct {
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	} `json:"author"`
}

type ApprovalRule struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	RuleType          string `json:"rule_type"`
	ApprovalsRequired int    `json:"approvals_required"`
	EligibleApprovers []User `json:"eligible_approvers"`
}

type User struct {
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	WebURL      string `json:"web_url"`
	Name        string `json:"name"`
	PublicEmail string `json:"public_email"`
}

type ApprovedBy struct {
	User       User      `json:"user"`
	ApprovedAt time.Time `json:"approved_at"`
}

type Approval struct {
	ApprovedBy []ApprovedBy `json:"approved_by"`
}

type Note struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	Resolvable bool      `json:"resolvable"`
	Resolved   bool      `json:"resolved"`
	Author     User      `json:"author"`
	ResolvedBy User      `json:"resolved_by"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	ResolvedAt time.Time `json:"resolved_at"`
}

type Discussion struct {
	ID             string `json:"id"`
	IndividualNote bool   `json:"individual_note"`
	Notes          []Note `json:"notes"`
}

type Commit struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
}

type MergeRequestInfo struct {
	Author      User      `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	Description string    `json:"description"`
	Reviewers   []User    `json:"reviewers"`
	WebURL      string    `json:"web_url"`
	Pipeline    struct {
		Status string `json:"status"`
	} `json:"pipeline"`
}
