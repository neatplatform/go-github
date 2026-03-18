package rest

import (
	"context"
	"fmt"
	"time"
)

// IssueService provides GitHub REST APIs for issues in a repository.
//
// See https://docs.github.com/en/rest/issues
type IssueService struct {
	client      *Client
	owner, repo string
}

type (
	// PullURLs is an object added to an issue representing a pull request.
	PullURLs struct {
		URL      string `json:"url"`
		HTMLURL  string `json:"html_url"`
		DiffURL  string `json:"diff_url"`
		PatchURL string `json:"patch_url"`
	}

	// Issue is a GitHub issue object.
	Issue struct {
		ID        int        `json:"id"`
		Number    int        `json:"number"`
		State     string     `json:"state"`
		Locked    bool       `json:"locked"`
		Title     string     `json:"title"`
		Body      string     `json:"body"`
		User      User       `json:"user"`
		Labels    []Label    `json:"labels"`
		Milestone *Milestone `json:"milestone"`
		URL       string     `json:"url"`
		HTMLURL   string     `json:"html_url"`
		LabelsURL string     `json:"labels_url"`
		PullURLs  *PullURLs  `json:"pull_request"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		ClosedAt  *time.Time `json:"closed_at"`
	}
)

// Event is a GitHub event object.
type Event struct {
	ID        int       `json:"id"`
	Event     string    `json:"event"`
	CommitID  string    `json:"commit_id"`
	Actor     User      `json:"actor"`
	URL       string    `json:"url"`
	CommitURL string    `json:"commit_url"`
	CreatedAt time.Time `json:"created_at"`
}

// IssuesFilter are used for fetching Issues.
type IssuesFilter struct {
	State string
	Since time.Time
}

// List retrieves all issues in the repository page by page.
//
// See https://docs.github.com/en/rest/issues/issues#list-repository-issues
func (s *IssueService) List(ctx context.Context, pageSize, pageNo int, filter IssuesFilter) ([]Issue, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/issues", s.owner, s.repo)
	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	if filter.State != "" {
		q.Add("state", filter.State)
	}
	if !filter.Since.IsZero() {
		q.Add("since", filter.Since.Format(time.RFC3339))
	}
	req.URL.RawQuery = q.Encode()

	issues := []Issue{}
	resp, err := s.client.Do(req, &issues)
	if err != nil {
		return nil, nil, err
	}

	return issues, resp, nil
}

// Events retrieves all events for an issue in the repository page by page.
//
// See https://docs.github.com/en/rest/issues/events#list-issue-events
func (s *IssueService) Events(ctx context.Context, number, pageSize, pageNo int) ([]Event, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/issues/%d/events", s.owner, s.repo, number)
	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	events := []Event{}

	resp, err := s.client.Do(req, &events)
	if err != nil {
		return nil, nil, err
	}

	return events, resp, nil
}
