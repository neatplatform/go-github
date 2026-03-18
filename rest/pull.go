package rest

import (
	"context"
	"fmt"
	"time"
)

// PullService provides GitHub REST APIs for pull requests in a repository.
//
// See https://docs.github.com/en/rest/pulls
type PullService struct {
	client      *Client
	owner, repo string
}

type (
	// CreatePullParams is used for creating a pull request.
	//
	// See https://docs.github.com/en/rest/pulls/pulls#create-a-pull-request
	CreatePullParams struct {
		Draft bool   `json:"draft"`
		Title string `json:"title"`
		Body  string `json:"body"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	}

	// UpdatePullParams is used for updating a pull request.
	//
	// See https://docs.github.com/en/rest/pulls/pulls#update-a-pull-request
	UpdatePullParams struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Base  string `json:"base"`
		State string `json:"state"` // Either open or closed
	}

	// PullBranch represents a base or head object in a Pull object.
	PullBranch struct {
		Label string     `json:"label"`
		Ref   string     `json:"ref"`
		SHA   string     `json:"sha"`
		User  User       `json:"user"`
		Repo  Repository `json:"repo"`
	}

	// Pull is a GitHub pull request object.
	Pull struct {
		ID             int        `json:"id"`
		Number         int        `json:"number"`
		State          string     `json:"state"`
		Draft          bool       `json:"draft"`
		Locked         bool       `json:"locked"`
		Title          string     `json:"title"`
		Body           string     `json:"body"`
		User           User       `json:"user"`
		Labels         []Label    `json:"labels"`
		Milestone      *Milestone `json:"milestone"`
		Base           PullBranch `json:"base"`
		Head           PullBranch `json:"head"`
		Merged         bool       `json:"merged"`
		Mergeable      *bool      `json:"mergeable"`
		Rebaseable     *bool      `json:"rebaseable"`
		MergedBy       *User      `json:"merged_by"`
		MergeCommitSHA string     `json:"merge_commit_sha"`
		URL            string     `json:"url"`
		HTMLURL        string     `json:"html_url"`
		DiffURL        string     `json:"diff_url"`
		PatchURL       string     `json:"patch_url"`
		IssueURL       string     `json:"issue_url"`
		CommitsURL     string     `json:"commits_url"`
		StatusesURL    string     `json:"statuses_url"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
		ClosedAt       *time.Time `json:"closed_at"`
		MergedAt       *time.Time `json:"merged_at"`
	}
)

// Get retrieves a pull request in the repository by its number.
//
// See https://docs.github.com/en/rest/pulls/pulls#get-a-pull-request
func (s *PullService) Get(ctx context.Context, number int) (*Pull, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/pulls/%d", s.owner, s.repo, number)
	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	pull := new(Pull)

	resp, err := s.client.Do(req, pull)
	if err != nil {
		return nil, nil, err
	}

	return pull, resp, nil
}

// PullsFilter are used for fetching Pulls.
//
// See https://docs.github.com/en/rest/pulls/pulls#list-pull-requests
type PullsFilter struct {
	State string
}

// List retrieves all pull requests in the repository page by page.
//
// See https://docs.github.com/en/rest/pulls/pulls#list-pull-requests
func (s *PullService) List(ctx context.Context, pageSize, pageNo int, filter PullsFilter) ([]Pull, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/pulls", s.owner, s.repo)
	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	if filter.State != "" {
		q.Add("state", filter.State)
	}
	req.URL.RawQuery = q.Encode()

	pulls := []Pull{}
	resp, err := s.client.Do(req, &pulls)
	if err != nil {
		return nil, nil, err
	}

	return pulls, resp, nil
}

// Create creates a new pull request in the repository.
//
// See https://docs.github.com/en/rest/pulls/pulls#create-a-pull-request
func (s *PullService) Create(ctx context.Context, params CreatePullParams) (*Pull, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/pulls", s.owner, s.repo)
	req, err := s.client.NewRequest(ctx, "POST", url, params)
	if err != nil {
		return nil, nil, err
	}

	pull := new(Pull)

	resp, err := s.client.Do(req, pull)
	if err != nil {
		return nil, nil, err
	}

	return pull, resp, nil
}

// Update updates a pull request in the repository.
//
// See https://docs.github.com/en/rest/pulls/pulls#update-a-pull-request
func (s *PullService) Update(ctx context.Context, number int, params UpdatePullParams) (*Pull, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/pulls/%d", s.owner, s.repo, number)
	req, err := s.client.NewRequest(ctx, "PATCH", url, params)
	if err != nil {
		return nil, nil, err
	}

	pull := new(Pull)

	resp, err := s.client.Do(req, pull)
	if err != nil {
		return nil, nil, err
	}

	return pull, resp, nil
}
