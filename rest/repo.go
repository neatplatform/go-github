package rest

import (
	"context"
	"fmt"
	"io"
	"time"
)

// RepoService provides GitHub REST APIs for a repository.
//
// See https://docs.github.com/en/rest/reference/repos
type RepoService struct {
	client      *Client
	owner, repo string

	// Services
	Pulls    *PullService
	Issues   *IssueService
	Releases *ReleaseService
}

// Repository is a GitHub repository object.
type Repository struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Description   string    `json:"description"`
	Topics        []string  `json:"topics"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	Archived      bool      `json:"archived"`
	Disabled      bool      `json:"disabled"`
	DefaultBranch string    `json:"default_branch"`
	Owner         User      `json:"owner"`
	URL           string    `json:"url"`
	HTMLURL       string    `json:"html_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PushedAt      time.Time `json:"pushed_at"`
}

// Permission represents a GitHub repository permission.
//
// See https://docs.github.com/en/github/setting-up-and-managing-organizations-and-teams/repository-permission-levels-for-an-organization
type Permission string

const (
	// PermissionNone does not allow anything.
	PermissionNone Permission = "none"
	// PermissionRead allows a contributor to view or discuss a project.
	PermissionRead Permission = "read"
	// PermissionTriage allows a contributor to manage issues and pull requests without write access.
	PermissionTriage Permission = "triage"
	// PermissionWrite allows a contributor to push to a project.
	PermissionWrite Permission = "write"
	// PermissionMaintain allows a contributor to manage a repository without access to sensitive or destructive actions.
	PermissionMaintain Permission = "maintain"
	// PermissionAdmin gives a contributor full access to a project, including sensitive and destructive actions.
	PermissionAdmin Permission = "admin"
)

type (
	// Hash is a GitHub hash object.
	Hash struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	}

	// Signature is a GitHub signature object.
	Signature struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Time  time.Time `json:"date"`
	}

	// RawCommit is a GitHub raw commit object.
	RawCommit struct {
		Message   string    `json:"message"`
		Author    Signature `json:"author"`
		Committer Signature `json:"committer"`
		Tree      Hash      `json:"tree"`
		URL       string    `json:"url"`
	}

	// Commit is a GitHub repository commit object.
	Commit struct {
		SHA       string    `json:"sha"`
		Commit    RawCommit `json:"commit"`
		Author    User      `json:"author"`
		Committer User      `json:"committer"`
		Parents   []Hash    `json:"parents"`
		URL       string    `json:"url"`
		HTMLURL   string    `json:"html_url"`
	}
)

// Branch is a GitHub branch object.
type Branch struct {
	Name      string `json:"name"`
	Protected bool   `json:"protected"`
	Commit    Commit `json:"commit"`
}

// Tag is a GitHib tag object.
type Tag struct {
	Name   string `json:"name"`
	Commit Hash   `json:"commit"`
}

// Label is a GitHub label object.
type Label struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
	URL         string `json:"url"`
}

// Milestone is a GitHub milestone object.
type Milestone struct {
	ID           int        `json:"id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      User       `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	DueOn        *time.Time `json:"due_on"`
	URL          string     `json:"url"`
	HTMLURL      string     `json:"html_url"`
	LabelsURL    string     `json:"labels_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
}

type (
	// ReleaseParams is used for creating or updating a GitHub release.
	ReleaseParams struct {
		Name       string `json:"name"`
		TagName    string `json:"tag_name"`
		Target     string `json:"target_commitish"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
		Body       string `json:"body"`
	}

	// Release is a GitHub release object.
	Release struct {
		ID          int            `json:"id"`
		Name        string         `json:"name"`
		TagName     string         `json:"tag_name"`
		Target      string         `json:"target_commitish"`
		Draft       bool           `json:"draft"`
		Prerelease  bool           `json:"prerelease"`
		Body        string         `json:"body"`
		URL         string         `json:"url"`
		HTMLURL     string         `json:"html_url"`
		AssetsURL   string         `json:"assets_url"`
		UploadURL   string         `json:"upload_url"`
		TarballURL  string         `json:"tarball_url"`
		ZipballURL  string         `json:"zipball_url"`
		CreatedAt   time.Time      `json:"created_at"`
		PublishedAt time.Time      `json:"published_at"`
		Author      User           `json:"author"`
		Assets      []ReleaseAsset `json:"assets"`
	}

	// ReleaseAsset is a Github release asset object.
	ReleaseAsset struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Label         string    `json:"label"`
		State         string    `json:"state"`
		ContentType   string    `json:"content_type"`
		Size          int       `json:"size"`
		DownloadCount int       `json:"download_count"`
		URL           string    `json:"url"`
		DownloadURL   string    `json:"browser_download_url"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Uploader      User      `json:"uploader"`
	}
)

// Get retrieves a repository by its name.
//
// See https://docs.github.com/rest/reference/repos#get-a-repository
func (s *RepoService) Get(ctx context.Context) (*Repository, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s", s.owner, s.repo)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	repository := new(Repository)

	resp, err := s.client.Do(req, repository)
	if err != nil {
		return nil, nil, err
	}

	return repository, resp, nil
}

// Permission returns the repository permission for a collaborator (user).
//
// See https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
func (s *RepoService) Permission(ctx context.Context, username string) (Permission, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/collaborators/%s/permission", s.owner, s.repo, username)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	body := new(struct {
		Permission Permission `json:"permission"`
		User       User       `json:"user"`
	})

	resp, err := s.client.Do(req, body)
	if err != nil {
		return "", nil, err
	}

	return body.Permission, resp, nil
}

// Commit retrieves a commit in the repository by its reference.
//
// See https://docs.github.com/rest/reference/repos#get-a-commit
func (s *RepoService) Commit(ctx context.Context, ref string) (*Commit, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/commits/%s", s.owner, s.repo, ref)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	commit := new(Commit)

	resp, err := s.client.Do(req, commit)
	if err != nil {
		return nil, nil, err
	}

	return commit, resp, nil
}

// Commits retrieves all commits in the repository page by page.
//
// See https://docs.github.com/rest/reference/repos#list-commits
func (s *RepoService) Commits(ctx context.Context, pageSize, pageNo int) ([]Commit, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/commits", s.owner, s.repo)
	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	commits := []Commit{}

	resp, err := s.client.Do(req, &commits)
	if err != nil {
		return nil, nil, err
	}

	return commits, resp, nil
}

// Branch retrieves a branch in the repository by its name.
//
// See https://docs.github.com/rest/reference/repos#get-a-branch
func (s *RepoService) Branch(ctx context.Context, name string) (*Branch, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/branches/%s", s.owner, s.repo, name)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	branch := new(Branch)

	resp, err := s.client.Do(req, branch)
	if err != nil {
		return nil, nil, err
	}

	return branch, resp, nil
}

// BranchProtection enables/disables a branch protection for administrator users.
//
// See https://docs.github.com/rest/reference/repos#set-admin-branch-protection
// See https://docs.github.com/rest/reference/repos#delete-admin-branch-protection
func (s *RepoService) BranchProtection(ctx context.Context, branch string, enabled bool) (*Response, error) {
	var method string
	if enabled {
		method = "POST"
	} else {
		method = "DELETE"
	}

	url := fmt.Sprintf("/repos/%s/%s/branches/%s/protection/enforce_admins", s.owner, s.repo, branch)
	req, err := s.client.NewRequest(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Tags retrieves all tags in the repository page by page.
//
// See https://docs.github.com/rest/reference/repos#list-repository-tags
func (s *RepoService) Tags(ctx context.Context, pageSize, pageNo int) ([]Tag, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/tags", s.owner, s.repo)
	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	tags := []Tag{}

	resp, err := s.client.Do(req, &tags)
	if err != nil {
		return nil, nil, err
	}

	return tags, resp, nil
}

// DownloadTarArchive downloads a repository archive in tar format.
func (s *RepoService) DownloadTarArchive(ctx context.Context, ref string, w io.Writer) (*Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/tarball/%s", s.owner, s.repo, ref)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DownloadZipArchive downloads a repository archive in zip format.
func (s *RepoService) DownloadZipArchive(ctx context.Context, ref string, w io.Writer) (*Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/zipball/%s", s.owner, s.repo, ref)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
