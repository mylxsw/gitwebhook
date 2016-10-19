package gitlab

import "time"

type GitLabObj struct {
	ObjectKind       string        `json:"object_kind"`
	Before           string        `json:"before"`
	After            string        `json:"after"`
	Ref              string        `json:"ref"`
	CheckoutSha      string        `json:"checkout_sha"`
	Message          string        `json:"message"`
	UserId           int           `json:"user_id"`
	Username         string        `json:"username"`
	UserEmail        string        `json:"user_email"`
	ProjectId        int           `json:"project_id"`
	Repository       GitRepository `json:"repository"`
	Commits          []GitCommit   `json:"commits"`
	TotalCommitCount int           `json:"total_commits_count"`

	User             GitUser             `json:"user"`
	ObjectAttributes GitObjectAttributes `json:"object_attributes"`
	MergeRequest     GitMergeRequest     `json:"merge_request"`
	Issue            GitIssue            `json:"issue"`
	Snippet          GitSnippet          `json:"snippet"`
}

type GitRepository struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitHttpUrl      string `json:"git_http_url"`
	GitSSHUrl       string `json:"git_ssh_url"`
	VisibilityLevel int    `json:"visibility_level"`
}

type GitCommit struct {
	Id        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Url       string    `json:"url"`
	Author    GitAuthor `json:"author"`
}

type GitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GitUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
}

type GitObjectAttributes struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	AssigneeId  int       `json:"assignee_id"`
	AuthorId    int       `json:"author_id"`
	ProjectId   int       `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Position    int       `json:"position"`
	BranchName  string    `json:"branch_name"`
	Description string    `json:"description"`
	MilestoneId string    `json:"milestone_id"`
	State       string    `json:"state"`
	Iid         int       `json:"iid"`
	Url         string    `json:"url"`
	Action      string    `json:"action"`

	Note         string           `json:"note"`
	NoteableType string           `json:"noteable_type"`
	Attachment   string           `json:"attachment"`
	LineCode     string           `json:"line_code"`
	CommitId     string           `json:"commit_id"`
	NoteableId   string           `json:"noteable_id"`
	System       string           `json:"system"`
	StDiff       GitStDiff        `json:"st_diff"`
	Source       GitRequestTarget `json:"source"`
	Target       GitRequestTarget `json:"target"`
	LastCommit   GitCommit        `json:"last_commit"`
}

type GitStDiff struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"new_path"`
	OldPath     string `json:"old_path"`
	AMode       string `json:"a_mode"`
	BMode       string `json:"b_mode"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
	DeletedFile bool   `json:"deleted_file"`
}

type GitMergeRequest struct {
	Id              int              `json:"id"`
	TargetBranch    string           `json:"target_branch"`
	SourceBranch    string           `json:"source_branch"`
	SourceProjectId int              `json:"source_project_id"`
	AuthorId        int              `json:"author_id"`
	AssigneeId      int              `json:"assignee_id"`
	Title           string           `json:"title"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	MilestoneId     int              `json:"milestone_id"`
	State           string           `json:"state"`
	MergeStatus     string           `json:"merge_status"`
	TargetProjectId int              `json:"target_project_id"`
	Iid             int              `json:"iid"`
	Description     string           `json:"description"`
	Position        int              `json:"position"`
	LockedAt        time.Time        `json:"locked_at"`
	Source          GitRequestTarget `json:"source"`
	Target          GitRequestTarget `json:"target"`
	LastCommit      GitCommit        `json:"last_commit"`
}

type GitRequestTarget struct {
	Name            string `json:"name"`
	SSHUrl          string `json:"ssh_url"`
	HttpUrl         string `json:"http_url"`
	Namespace       string `json:"namespace"`
	VisibilityLevel int    `json:"visibility_level"`
}

type GitIssue struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	AssigneeId  int       `json:"assignee_id"`
	AuthorId    int       `json:"author_id"`
	ProjectId   int       `json:"project_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Position    int       `json:"position"`
	BranchName  string    `json:"branch_name"`
	Description string    `json:"description"`
	MilestoneId string    `json:"milestone_id"`
	State       string    `json:"state"`
	Iid         int       `json:"iid"`
}

type GitSnippet struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	AuthorId        int       `json:"author_id"`
	ProjectId       int       `json:"project_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	FileName        string    `json:"file_name"`
	ExpiresAt       time.Time `json:"expires_at"`
	Type            string    `json:"type"`
	VisibilityLevel int       `json:"visibility_level"`
}

const (
	ObjectKindPush    = "push"
	ObjectKindTagPush = "tag_push"
	ObjectKindIssue   = "issue"
	ObjectKindNote    = "note"
	ObjectKindMerge   = "merge_request"

	NoteableTypeCommit  = "Commit"
	NoteableTypeMerge   = "MergeRequest"
	NoteableTypeIssue   = "Issue"
	NoteableTypeSnippet = "Snippet"
)
