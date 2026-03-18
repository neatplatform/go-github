package rest

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pullBody = `{
		"id": 1,
		"url": "https://api.github.com/repos/octocat/Hello-World/pulls/1001",
		"html_url": "https://github.com/octocat/Hello-World/pull/1001",
		"number": 1001,
		"state": "open",
		"locked": false,
		"draft": false,
		"title":  "Add new feature",
		"body": "This is an awesome feature!",
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
				"name": "feature"
			}
		],
		"milestone": {
			"id": 3000,
			"number": 1,
			"state": "open",
			"title": "v1.0"
		},
		"created_at":  "2020-10-15T15:00:00Z",
		"updated_at": "2020-10-22T22:00:00Z",
		"head": {
			"label": "octocat:feature-branch",
			"ref": "feature-branch",
			"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		},
		"base": {
			"label": "octocat:main",
			"ref": "main",
			"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
		}
	}`

	pullsBody = `[
		{
			"id": 1,
			"url": "https://api.github.com/repos/octocat/Hello-World/pulls/1001",
			"html_url": "https://github.com/octocat/Hello-World/pull/1001",
			"number": 1001,
			"state": "open",
			"locked": false,
			"draft": false,
			"title": "Add new feature",
			"body": "This is an awesome feature!",
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
					"name": "feature"
				}
			],
			"milestone": {
				"id": 3000,
				"number": 1,
				"state": "open",
				"title": "v1.0"
			},
			"created_at":  "2020-10-15T15:00:00Z",
			"updated_at": "2020-10-22T22:00:00Z",
			"head": {
				"label": "octocat:feature-branch",
				"ref": "feature-branch",
				"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
			},
			"base": {
				"label": "octocat:main",
				"ref": "main",
				"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
			}
		}
	]`
)

var (
	pull = Pull{
		ID:     1,
		Number: 1001,
		State:  "open",
		Draft:  false,
		Locked: false,
		Title:  "Add new feature",
		Body:   "This is an awesome feature!",
		User: User{
			ID:      1,
			Login:   "octocat",
			Type:    "User",
			URL:     "https://api.github.com/users/octocat",
			HTMLURL: "https://github.com/octocat",
		},
		Labels: []Label{
			{
				ID:   2000,
				Name: "feature",
			},
		},
		Milestone: &Milestone{
			ID:     3000,
			Number: 1,
			State:  "open",
			Title:  "v1.0",
		},
		Base: PullBranch{
			Label: "octocat:main",
			Ref:   "main",
			SHA:   "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		},
		Head: PullBranch{
			Label: "octocat:feature-branch",
			Ref:   "feature-branch",
			SHA:   "6dcb09b5b57875f334f61aebed695e2e4193db5e",
		},
		Merged:    false,
		URL:       "https://api.github.com/repos/octocat/Hello-World/pulls/1001",
		HTMLURL:   "https://github.com/octocat/Hello-World/pull/1001",
		CreatedAt: parseGitHubTime("2020-10-15T15:00:00Z"),
		UpdatedAt: parseGitHubTime("2020-10-22T22:00:00Z"),
	}
)

func TestPullService_Get(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *PullService
		ctx              context.Context
		number           int
		expectedPull     *Pull
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			number:        1001,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls/1001", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1001,
			expectedError: `GET /repos/octocat/Hello-World/pulls/1001: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls/1001", 200, http.Header{}, `{`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1001,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls/1001", 200, header, pullBody},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:          context.Background(),
			number:       1001,
			expectedPull: &pull,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			pull, resp, err := tc.s.Get(tc.ctx, tc.number)

			if tc.expectedError != "" {
				assert.Nil(t, pull)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPull, pull)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestPullService_List(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *PullService
		ctx              context.Context
		pageSize         int
		pageNo           int
		params           PullsFilter
		expectedPulls    []Pull
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      nil,
			pageSize: 10,
			pageNo:   1,
			params: PullsFilter{
				State: "closed",
			},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: PullsFilter{
				State: "closed",
			},
			expectedError: `GET /repos/octocat/Hello-World/pulls: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls", 200, http.Header{}, `[`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: PullsFilter{
				State: "closed",
			},
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/repos/octocat/Hello-World/pulls", 200, header, pullsBody},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:      context.Background(),
			pageSize: 10,
			pageNo:   1,
			params: PullsFilter{
				State: "closed",
			},
			expectedPulls: []Pull{pull},
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

			pulls, resp, err := tc.s.List(tc.ctx, tc.pageSize, tc.pageNo, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, pulls)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPulls, pulls)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestPullService_Create(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	params := CreatePullParams{
		Draft: false,
		Title: "Add new feature",
		Body:  "This is an awesome feature!",
		Head:  "octocat:feature-branch",
		Base:  "octocat:main",
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *PullService
		ctx              context.Context
		params           CreatePullParams
		expectedPull     *Pull
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			params:        params,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/pulls", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			params:        params,
			expectedError: `POST /repos/octocat/Hello-World/pulls: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/pulls", 201, http.Header{}, `{`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			params:        params,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"POST", "/repos/octocat/Hello-World/pulls", 201, header, pullBody},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:          context.Background(),
			params:       params,
			expectedPull: &pull,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			pull, resp, err := tc.s.Create(tc.ctx, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, pull)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPull, pull)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestPullService_Update(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicAPIURL,
	}

	params := UpdatePullParams{
		Title: "[CLOSED] Add new feature",
		Body:  "[CLOSED] This is an awesome feature!",
		Base:  "octocat:main",
		State: "closed",
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *PullService
		ctx              context.Context
		number           int
		params           UpdatePullParams
		expectedPull     *Pull
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           nil,
			number:        1,
			params:        params,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/pulls/1", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1,
			params:        params,
			expectedError: `PATCH /repos/octocat/Hello-World/pulls/1: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/pulls/1", 200, http.Header{}, `{`},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:           context.Background(),
			number:        1,
			params:        params,
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"PATCH", "/repos/octocat/Hello-World/pulls/1", 200, header, pullBody},
			},
			s: &PullService{
				client: c,
				owner:  "octocat",
				repo:   "Hello-World",
			},
			ctx:          context.Background(),
			number:       1,
			params:       params,
			expectedPull: &pull,
			expectedResponse: &Response{
				Rate: expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			tc.s.client.apiURL, _ = url.Parse(ts.URL)

			pull, resp, err := tc.s.Update(tc.ctx, tc.number, tc.params)

			if tc.expectedError != "" {
				assert.Nil(t, pull)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPull, pull)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}
