package github

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	issuesBody = `[
		{
			"id": 2,
			"url": "https://api.github.com/repos/octocat/Hello-World/issues/1002",
			"html_url": "https://github.com/octocat/Hello-World/pull/1002",
			"number": 1002,
			"state": "closed",
			"title": "Fixed a bug",
			"body": "I made this to work as expected!",
			"user": {
				"login": "octodog",
				"id": 2,
				"url": "https://api.github.com/users/octodog",
				"html_url": "https://github.com/octodog",
				"type": "User"
			},
			"labels": [
				{
					"id": 2000,
					"name": "bug",
					"default": true
				}
			],
			"milestone": {
				"id": 3000,
				"number": 1,
				"state": "open",
				"title": "v1.0"
			},
			"locked": false,
			"pull_request": {
				"url": "https://api.github.com/repos/octocat/Hello-World/pulls/1002"
			},
			"closed_at": "2020-10-20T20:00:00Z",
			"created_at": "2020-10-15T15:00:00Z",
			"updated_at": "2020-10-22T22:00:00Z"
		},
		{
			"id": 1,
			"url": "https://api.github.com/repos/octocat/Hello-World/issues/1001",
			"html_url": "https://github.com/octocat/Hello-World/issues/1001",
			"number": 1001,
			"state": "open",
			"title": "Found a bug",
			"body": "This is not working as expected!",
			"user": {
				"login": "octocat",
				"id": 1,
				"url": "https://api.github.com/users/octocat",
				"html_url": "https://github.com/octocat",
				"type": "User"
			},
			"labels": [
				{
					"id": 2000,
					"name": "bug",
					"default": true
				}
			],
			"milestone": {
				"id": 3000,
				"number": 1,
				"state": "open",
				"title": "v1.0"
			},
			"locked": true,
			"pull_request": null,
			"closed_at": null,
			"created_at": "2020-10-10T10:00:00Z",
			"updated_at": "2020-10-20T20:00:00Z"
		}
	]`

	eventsBody = `[
		{
			"id": 2,
			"actor": {
				"login": "octofox",
				"id": 3,
				"url": "https://api.github.com/users/octofox",
				"html_url": "https://github.com/octofox",
				"type": "User"
			},
			"event": "merged",
			"commit_id": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			"created_at": "2020-10-20T20:00:00Z"
		},
		{
			"id": 1,
			"actor": {
				"login": "octocat",
				"id": 1,
				"url": "https://api.github.com/users/octocat",
				"html_url": "https://github.com/octocat",
				"type": "User"
			},
			"event": "closed",
			"commit_id": null,
			"created_at": "2020-10-20T20:00:00Z"
		}
	]`
)

var (
	issue1 = Issue{
		ID:     1,
		Number: 1001,
		State:  "open",
		Locked: true,
		Title:  "Found a bug",
		Body:   "This is not working as expected!",
		User: User{
			ID:      1,
			Login:   "octocat",
			Type:    "User",
			URL:     "https://api.github.com/users/octocat",
			HTMLURL: "https://github.com/octocat",
		},
		Labels: []Label{
			{
				ID:      2000,
				Name:    "bug",
				Default: true,
			},
		},
		Milestone: &Milestone{
			ID:     3000,
			Number: 1,
			State:  "open",
			Title:  "v1.0",
		},
		URL:       "https://api.github.com/repos/octocat/Hello-World/issues/1001",
		HTMLURL:   "https://github.com/octocat/Hello-World/issues/1001",
		CreatedAt: parseGitHubTime("2020-10-10T10:00:00Z"),
		UpdatedAt: parseGitHubTime("2020-10-20T20:00:00Z"),
		ClosedAt:  nil,
	}

	issue2 = Issue{
		ID:     2,
		Number: 1002,
		State:  "closed",
		Locked: false,
		Title:  "Fixed a bug",
		Body:   "I made this to work as expected!",
		User: User{
			ID:      2,
			Login:   "octodog",
			Type:    "User",
			URL:     "https://api.github.com/users/octodog",
			HTMLURL: "https://github.com/octodog",
		},
		Labels: []Label{
			{
				ID:      2000,
				Name:    "bug",
				Default: true,
			},
		},
		Milestone: &Milestone{
			ID:     3000,
			Number: 1,
			State:  "open",
			Title:  "v1.0",
		},
		URL:     "https://api.github.com/repos/octocat/Hello-World/issues/1002",
		HTMLURL: "https://github.com/octocat/Hello-World/pull/1002",
		PullURLs: &PullURLs{
			URL: "https://api.github.com/repos/octocat/Hello-World/pulls/1002",
		},
		CreatedAt: parseGitHubTime("2020-10-15T15:00:00Z"),
		UpdatedAt: parseGitHubTime("2020-10-22T22:00:00Z"),
		ClosedAt:  parseGitHubTimePtr("2020-10-20T20:00:00Z"),
	}

	event1 = Event{
		ID:       1,
		Event:    "closed",
		CommitID: "",
		Actor: User{
			ID:      1,
			Login:   "octocat",
			Type:    "User",
			URL:     "https://api.github.com/users/octocat",
			HTMLURL: "https://github.com/octocat",
		},
		CreatedAt: parseGitHubTime("2020-10-20T20:00:00Z"),
	}

	event2 = Event{
		ID:       2,
		Event:    "merged",
		CommitID: "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		Actor: User{
			ID:      3,
			Login:   "octofox",
			Type:    "User",
			URL:     "https://api.github.com/users/octofox",
			HTMLURL: "https://github.com/octofox",
		},
		CreatedAt: parseGitHubTime("2020-10-20T20:00:00Z"),
	}
)

func TestIssueService_List(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rateLimits: map[Category]Rate{},
		apiURL:     publicAPIURL,
	}

	since, _ := time.Parse(time.RFC3339, "2020-10-20T22:30:00-04:00")

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *IssueService
		ctx              context.Context
		pageSize         int
		pageNo           int
		params           IssuesFilter
		expectedIssues   []Issue
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      nil,
			pageSize: 10,
			pageNo:   1,
			params: IssuesFilter{
				State: "closed",
				Since: since,
			},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "StatusUnauthorized",
			mockResponses: []MockResponse{
				{
					"GET", "/repos/octocat/Hello-World/issues", 401, http.Header{}, `{
						"status": "401",
						"message": "Bad credentials"
					}`,
				},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: IssuesFilter{
				State: "closed",
				Since: since,
			},
			expectedError: `GET /repos/octocat/Hello-World/issues: [401] Bad credentials`,
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/issues", 200, http.Header{}, `[`},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: IssuesFilter{
				State: "closed",
				Since: since,
			},
			expectedError: `failed to decode response body: unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/issues", 200, header, issuesBody},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: IssuesFilter{
				State: "closed",
				Since: since,
			},
			expectedIssues: []Issue{issue2, issue1},
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			issues, resp, err := tc.s.List(tc.ctx, tc.pageSize, tc.pageNo, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, issues)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedIssues, issues)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestIssueService_Events(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rateLimits: map[Category]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *IssueService
		ctx              context.Context
		number           int
		pageSize         int
		pageNo           int
		expectedEvents   []Event
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			number:        1001,
			pageSize:      10,
			pageNo:        1,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "StatusUnauthorized",
			mockResponses: []MockResponse{
				{
					"GET", "/repos/octocat/Hello-World/issues/1001/events", 401, http.Header{}, `{
						"status": "401",
						"message": "Bad credentials"
					}`,
				},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1001,
			pageSize:      10,
			pageNo:        1,
			expectedError: `GET /repos/octocat/Hello-World/issues/1001/events: [401] Bad credentials`,
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/issues/1001/events", 200, http.Header{}, `[`},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1001,
			pageSize:      10,
			pageNo:        1,
			expectedError: `failed to decode response body: unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/issues/1001/events", 200, header, eventsBody},
			},
			s: &IssueService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:            context.Background(),
			number:         1001,
			pageSize:       10,
			pageNo:         1,
			expectedEvents: []Event{event2, event1},
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			events, resp, err := tc.s.Events(tc.ctx, tc.number, tc.pageSize, tc.pageNo)

			if tc.expectedError != "" {
				assert.Nil(t, events)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedEvents, events)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}
