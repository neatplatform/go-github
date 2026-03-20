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
	searchUsersBody = `{
		"total_count": 40,
		"incomplete_results": false,
		"items": [
			{
				"login": "octocat",
				"id": 1,
				"url": "https://api.github.com/users/octocat",
				"html_url": "https://github.com/octocat",
				"type": "User",
				"site_admin": false,
				"name": "The Octocat",
				"email": "octocat@github.com"
			}
		]
	}`

	searchReposBody = `{
		"total_count": 40,
		"incomplete_results": false,
		"items": [
			{
				"id": 1296269,
				"name": "Hello-World",
				"full_name": "octocat/Hello-World",
				"owner": {
					"login": "octocat",
					"id": 1,
					"type": "User"
				},
				"private": false,
				"description": "This your first repo!",
				"fork": false,
				"default_branch": "main",
				"topics": [
					"octocat",
					"api"
				],
				"archived": false,
				"disabled": false,
				"visibility": "public",
				"pushed_at": "2020-10-31T14:00:00Z",
				"created_at": "2020-01-20T09:00:00Z",
				"updated_at": "2020-10-31T14:00:00Z"
			}
		]
	}`

	searchIssuesBody = `{
		"total_count": 40,
		"incomplete_results": false,
		"items": [
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
			},
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
			}
		]
	}`
)

var (
	searchUsersResult = SearchUsersResult{
		TotalCount:        40,
		IncompleteResults: false,
		Items: []User{
			User{
				ID:      1,
				Login:   "octocat",
				Type:    "User",
				Email:   "octocat@github.com",
				Name:    "The Octocat",
				URL:     "https://api.github.com/users/octocat",
				HTMLURL: "https://github.com/octocat",
			},
		},
	}

	searchReposResult = SearchReposResult{
		TotalCount:        40,
		IncompleteResults: false,
		Items: []Repository{
			{
				ID:            1296269,
				Name:          "Hello-World",
				FullName:      "octocat/Hello-World",
				Description:   "This your first repo!",
				Topics:        []string{"octocat", "api"},
				Private:       false,
				Fork:          false,
				Archived:      false,
				Disabled:      false,
				DefaultBranch: "main",
				Owner: User{
					ID:    1,
					Login: "octocat",
					Type:  "User",
				},
				CreatedAt: parseGitHubTime("2020-01-20T09:00:00Z"),
				UpdatedAt: parseGitHubTime("2020-10-31T14:00:00Z"),
				PushedAt:  parseGitHubTime("2020-10-31T14:00:00Z"),
			},
		},
	}

	searchIssuesResult = SearchIssuesResult{
		TotalCount:        40,
		IncompleteResults: false,
		Items: []Issue{
			{
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
			},
			{
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
			},
		},
	}
)

func TestQualifier(t *testing.T) {
	t.Run("QualifierUser", func(t *testing.T) {
		assert.Equal(t, Qualifier("user:octocat"), QualifierUser("octocat"))
	})

	t.Run("QualifierOrg", func(t *testing.T) {
		assert.Equal(t, Qualifier("org:example"), QualifierOrg("example"))
	})

	t.Run("QualifierRepo", func(t *testing.T) {
		assert.Equal(t, Qualifier("repo:octocat/Hello-World"), QualifierRepo("octocat", "Hello-World"))
	})

	t.Run("QualifierAuthor", func(t *testing.T) {
		assert.Equal(t, Qualifier("author:octocat"), QualifierAuthor("octocat"))
	})

	t.Run("QualifierAuthorApp", func(t *testing.T) {
		assert.Equal(t, Qualifier("author:app/bot"), QualifierAuthorApp("bot"))
	})

	t.Run("QualifierAssignee", func(t *testing.T) {
		assert.Equal(t, Qualifier("assignee:octocat"), QualifierAssignee("octocat"))
	})

	t.Run("QualifierLabel", func(t *testing.T) {
		assert.Equal(t, Qualifier(`label:"bug"`), QualifierLabel("bug"))
	})

	t.Run("QualifierMilestone", func(t *testing.T) {
		assert.Equal(t, Qualifier(`milestone:"v1"`), QualifierMilestone("v1"))
	})

	t.Run("QualifierProject", func(t *testing.T) {
		assert.Equal(t, Qualifier(`project:"10"`), QualifierProject("10"))
	})

	t.Run("QualifierRepoProject", func(t *testing.T) {
		assert.Equal(t, Qualifier("project:octocat/Hello-World/1"), QualifierRepoProject("octocat", "Hello-World", "1"))
	})

	t.Run("QualifierHead", func(t *testing.T) {
		assert.Equal(t, Qualifier("head:feature"), QualifierHead("feature"))
	})

	t.Run("QualifierBase", func(t *testing.T) {
		assert.Equal(t, Qualifier("base:main"), QualifierBase("main"))
	})

	t.Run("QualifierLanguage", func(t *testing.T) {
		assert.Equal(t, Qualifier("language:go"), QualifierLanguage("go"))
	})

	t.Run("QualifierTopic", func(t *testing.T) {
		assert.Equal(t, Qualifier("topic:terraform"), QualifierTopic("terraform"))
	})

	t.Run("QualifierCreatedOn", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:2021-10-24"), QualifierCreatedOn(tm))
	})

	t.Run("QualifierCreatedAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:>2021-10-24"), QualifierCreatedAfter(tm))
	})

	t.Run("QualifierCreatedOnOrAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:>=2021-10-24"), QualifierCreatedOnOrAfter(tm))
	})

	t.Run("QualifierCreatedBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:<2021-10-24"), QualifierCreatedBefore(tm))
	})

	t.Run("QualifierCreatedOnOrBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:<=2021-10-24"), QualifierCreatedOnOrBefore(tm))
	})

	t.Run("QualifierCreatedBetween", func(t *testing.T) {
		from, _ := time.Parse(time.RFC3339, "2021-10-22T22:10:40-04:00")
		to, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("created:2021-10-22..2021-10-24"), QualifierCreatedBetween(from, to))
	})

	t.Run("QualifierUpdatedOn", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:2021-10-24"), QualifierUpdatedOn(tm))
	})

	t.Run("QualifierUpdatedAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:>2021-10-24"), QualifierUpdatedAfter(tm))
	})

	t.Run("QualifierUpdatedOnOrAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:>=2021-10-24"), QualifierUpdatedOnOrAfter(tm))
	})

	t.Run("QualifierUpdatedBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:<2021-10-24"), QualifierUpdatedBefore(tm))
	})

	t.Run("QualifierUpdatedOnOrBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:<=2021-10-24"), QualifierUpdatedOnOrBefore(tm))
	})

	t.Run("QualifierUpdatedBetween", func(t *testing.T) {
		from, _ := time.Parse(time.RFC3339, "2021-10-22T22:10:40-04:00")
		to, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("updated:2021-10-22..2021-10-24"), QualifierUpdatedBetween(from, to))
	})

	t.Run("QualifierClosedOn", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:2021-10-24"), QualifierClosedOn(tm))
	})

	t.Run("QualifierClosedAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:>2021-10-24"), QualifierClosedAfter(tm))
	})

	t.Run("QualifierClosedOnOrAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:>=2021-10-24"), QualifierClosedOnOrAfter(tm))
	})

	t.Run("QualifierClosedBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:<2021-10-24"), QualifierClosedBefore(tm))
	})

	t.Run("QualifierClosedOnOrBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:<=2021-10-24"), QualifierClosedOnOrBefore(tm))
	})

	t.Run("QualifierClosedBetween", func(t *testing.T) {
		from, _ := time.Parse(time.RFC3339, "2021-10-22T22:10:40-04:00")
		to, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("closed:2021-10-22..2021-10-24"), QualifierClosedBetween(from, to))
	})

	t.Run("QualifierMergedOn", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:2021-10-24"), QualifierMergedOn(tm))
	})

	t.Run("QualifierMergedAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:>2021-10-24"), QualifierMergedAfter(tm))
	})

	t.Run("QualifierMergedOnOrAfter", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:>=2021-10-24"), QualifierMergedOnOrAfter(tm))
	})

	t.Run("QualifierMergedBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:<2021-10-24"), QualifierMergedBefore(tm))
	})

	t.Run("QualifierMergedOnOrBefore", func(t *testing.T) {
		tm, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:<=2021-10-24"), QualifierMergedOnOrBefore(tm))
	})

	t.Run("QualifierMergedBetween", func(t *testing.T) {
		from, _ := time.Parse(time.RFC3339, "2021-10-22T22:10:40-04:00")
		to, _ := time.Parse(time.RFC3339, "2021-10-24T16:12:52-04:00")
		assert.Equal(t, Qualifier("merged:2021-10-22..2021-10-24"), QualifierMergedBetween(from, to))
	})
}

func TestSearchQuery(t *testing.T) {
	tests := []struct {
		name              string
		includeKeywords   []string
		excludeKeywords   []string
		includeQualifiers []Qualifier
		excludeQualifiers []Qualifier
		expectedString    string
	}{
		{
			name:              "Empty",
			includeKeywords:   []string{},
			excludeKeywords:   []string{},
			includeQualifiers: []Qualifier{},
			excludeQualifiers: []Qualifier{},
			expectedString:    "",
		},
		{
			name:              "OK",
			includeKeywords:   []string{"Implement"},
			excludeKeywords:   []string{"WIP"},
			includeQualifiers: []Qualifier{QualifierIsPR, QualifierIsOpen},
			excludeQualifiers: []Qualifier{QualifierStatusFailure},
			expectedString:    `"Implement" NOT "WIP" is:pr is:open -status:failure`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			q := new(SearchQuery)
			q.IncludeKeywords(tc.includeKeywords...)
			q.ExcludeKeywords(tc.excludeKeywords...)
			q.IncludeQualifiers(tc.includeQualifiers...)
			q.ExcludeQualifiers(tc.excludeQualifiers...)

			assert.Equal(t, tc.expectedString, q.String())
		})
	}
}

func TestSearchUsers(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rateLimits: map[Category]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *SearchService
		ctx              context.Context
		pageSize         int
		pageNo           int
		sort             SearchResultSort
		order            SearchResultOrder
		query            SearchQuery
		expectedResult   *SearchUsersResult
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s:             &SearchService{client: c},
			ctx:           nil,
			pageSize:      10,
			pageNo:        1,
			sort:          SortByJoined,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "StatusBadRequest",
			mockResponses: []MockResponse{
				{
					"GET", "/search/users", 400, http.Header{}, `{
						"status": "400",
						"message": "Invalid query"
					}`,
				},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByJoined,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `GET /search/users: [400] Invalid query`,
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/search/users", 200, http.Header{}, `{`},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByJoined,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `failed to decode response body: unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/search/users", 200, header, searchUsersBody},
			},
			s:              &SearchService{client: c},
			ctx:            context.Background(),
			pageSize:       10,
			pageNo:         1,
			sort:           SortByJoined,
			order:          DescOrder,
			query:          SearchQuery{},
			expectedResult: &searchUsersResult,
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

			result, resp, err := tc.s.SearchUsers(tc.ctx, tc.pageSize, tc.pageNo, tc.sort, tc.order, tc.query)

			if tc.expectedError != "" {
				assert.Nil(t, result)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestSearchRepos(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rateLimits: map[Category]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *SearchService
		ctx              context.Context
		pageSize         int
		pageNo           int
		sort             SearchResultSort
		order            SearchResultOrder
		query            SearchQuery
		expectedResult   *SearchReposResult
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s:             &SearchService{client: c},
			ctx:           nil,
			pageSize:      10,
			pageNo:        1,
			sort:          SortByStars,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "StatusBadRequest",
			mockResponses: []MockResponse{
				{
					"GET", "/search/repositories", 400, http.Header{}, `{
						"status": "400",
						"message": "Invalid query"
					}`,
				},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByStars,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `GET /search/repositories: [400] Invalid query`,
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/search/repositories", 200, http.Header{}, `{`},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByStars,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `failed to decode response body: unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/search/repositories", 200, header, searchReposBody},
			},
			s:              &SearchService{client: c},
			ctx:            context.Background(),
			pageSize:       10,
			pageNo:         1,
			sort:           SortByStars,
			order:          DescOrder,
			query:          SearchQuery{},
			expectedResult: &searchReposResult,
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

			result, resp, err := tc.s.SearchRepos(tc.ctx, tc.pageSize, tc.pageNo, tc.sort, tc.order, tc.query)

			if tc.expectedError != "" {
				assert.Nil(t, result)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestSearchIssues(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rateLimits: map[Category]Rate{},
		apiURL:     publicAPIURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *SearchService
		ctx              context.Context
		pageSize         int
		pageNo           int
		sort             SearchResultSort
		order            SearchResultOrder
		query            SearchQuery
		expectedResult   *SearchIssuesResult
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s:             &SearchService{client: c},
			ctx:           nil,
			pageSize:      10,
			pageNo:        1,
			sort:          SortByCreated,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "StatusBadRequest",
			mockResponses: []MockResponse{
				{
					"GET", "/search/issues", 400, http.Header{}, `{
						"status": "400",
						"message": "Invalid query"
					}`,
				},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByCreated,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `GET /search/issues: [400] Invalid query`,
		},
		{
			name: "InvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/search/issues", 200, http.Header{}, `{`},
			},
			s:             &SearchService{client: c},
			ctx:           context.Background(),
			pageSize:      10,
			pageNo:        1,
			sort:          SortByCreated,
			order:         DescOrder,
			query:         SearchQuery{},
			expectedError: `failed to decode response body: unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/search/issues", 200, header, searchIssuesBody},
			},
			s:              &SearchService{client: c},
			ctx:            context.Background(),
			pageSize:       10,
			pageNo:         1,
			sort:           SortByCreated,
			order:          DescOrder,
			query:          SearchQuery{},
			expectedResult: &searchIssuesResult,
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

			result, resp, err := tc.s.SearchIssues(tc.ctx, tc.pageSize, tc.pageNo, tc.sort, tc.order, tc.query)

			if tc.expectedError != "" {
				assert.Nil(t, result)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}
