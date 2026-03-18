package rest

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	header = http.Header{
		headerLink:          {`<https://api.github.com/repositories/100/issues?page=2>; rel="prev", <https://api.github.com/repositories/100/issues?page=4>; rel="next", <https://api.github.com/repositories/100/issues?page=6>; rel="last", <https://api.github.com/repositories/100/issues?page=1>; rel="first"`},
		headerRateLimit:     {"5000"},
		headerRateUsed:      {"10"},
		headerRateRemaining: {"4990"},
		headerRateReset:     {"1605083281"},
	}

	expectedPages = Pages{
		First: 1,
		Prev:  2,
		Next:  4,
		Last:  6,
	}

	expectedRate = Rate{
		Limit:     5000,
		Used:      10,
		Remaining: 4990,
		Reset:     Epoch(1605083281),
	}
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		authToken string
	}{
		{
			name:      "OK",
			authToken: "abcdefghijklmnopqrstuvwxyz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewClient(tc.authToken)

			assert.NotNil(t, c)
			assert.NotNil(t, c.httpClient)
			assert.NotNil(t, c.rates)
			assert.NotNil(t, c.apiURL)
			assert.NotNil(t, c.uploadURL)
			assert.NotNil(t, c.downloadURL)
			assert.Equal(t, tc.authToken, c.authToken)
			assert.NotNil(t, c.Users)
			assert.NotNil(t, c.Search)
		})
	}
}

func TestNewEnterpriseClient(t *testing.T) {
	tests := []struct {
		name          string
		apiURL        string
		uploadURL     string
		downloadURL   string
		authToken     string
		expectedError string
	}{
		{
			name:          "InvalidAPIURL",
			apiURL:        ":invalid",
			uploadURL:     "",
			authToken:     "auth-token",
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "InvalidUploadURL",
			apiURL:        "https://api.github.internal.com",
			uploadURL:     ":invalid",
			authToken:     "auth-token",
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "InvalidDownloadURL",
			apiURL:        "https://api.github.internal.com",
			uploadURL:     "https://uploads.github.internal.com",
			downloadURL:   ":invalid",
			authToken:     "auth-token",
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "Success",
			apiURL:        "https://api.github.internal.com",
			uploadURL:     "https://uploads.github.internal.com",
			downloadURL:   "https://github.internal.com",
			authToken:     "auth-token",
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewEnterpriseClient(tc.apiURL, tc.uploadURL, tc.downloadURL, tc.authToken)

			if tc.expectedError != "" {
				assert.Nil(t, c)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, c)
				assert.NotNil(t, c.httpClient)
				assert.NotNil(t, c.rates)
				assert.NotNil(t, c.apiURL)
				assert.NotNil(t, c.uploadURL)
				assert.NotNil(t, c.downloadURL)
				assert.Equal(t, tc.authToken, c.authToken)
				assert.NotNil(t, c.Users)
				assert.NotNil(t, c.Search)
			}
		})
	}
}

func TestClient_NewRequest(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		method        string
		url           string
		body          interface{}
		expectedError string
	}{
		{
			name:          "InvalidURL",
			ctx:           context.Background(),
			method:        "GET",
			url:           ":invalid",
			body:          nil,
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "InvalidBody",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			body:          make(chan int),
			expectedError: `json: unsupported type: chan int`,
		},
		{
			name:          "NilContext",
			ctx:           nil,
			method:        "GET",
			url:           "/user",
			body:          "request body",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "Success_Writer",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			body:          strings.NewReader("content"),
			expectedError: ``,
		},
		{
			name:          "Success_Struct",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			body:          new(struct{}),
			expectedError: ``,
		},
		{
			name:          "Success_Map",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			body:          make(map[string]interface{}),
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{
				apiURL:    publicAPIURL,
				authToken: "abcdefghijklmnopqrstuvwxyz",
			}

			req, err := c.NewRequest(tc.ctx, tc.method, tc.url, tc.body)

			if tc.expectedError != "" {
				assert.Nil(t, req)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.NotEmpty(t, req.Header.Get(headerUserAgent))
				assert.NotEmpty(t, req.Header.Get(headerAccept))
				assert.NotEmpty(t, req.Header.Get(headerAuth))
			}
		})
	}
}

func TestClient_NewPageRequest(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		method        string
		url           string
		pageSize      int
		pageNo        int
		body          interface{}
		expectedError string
	}{
		{
			name:          "NilContext",
			ctx:           nil,
			method:        "GET",
			url:           "/user",
			pageSize:      20,
			pageNo:        2,
			body:          "request body",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "Success_Writer",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			pageSize:      20,
			pageNo:        2,
			body:          strings.NewReader("content"),
			expectedError: ``,
		},
		{
			name:          "Success_Struct",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			pageSize:      20,
			pageNo:        2,
			body:          new(struct{}),
			expectedError: ``,
		},
		{
			name:          "Success_Map",
			ctx:           context.Background(),
			method:        "GET",
			url:           "/user",
			pageSize:      20,
			pageNo:        2,
			body:          make(map[string]interface{}),
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{
				apiURL:    publicAPIURL,
				authToken: "abcdefghijklmnopqrstuvwxyz",
			}

			req, err := c.NewPageRequest(tc.ctx, tc.method, tc.url, tc.pageSize, tc.pageNo, tc.body)

			if tc.expectedError != "" {
				assert.Nil(t, req)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.NotEmpty(t, req.Header.Get(headerUserAgent))
				assert.NotEmpty(t, req.Header.Get(headerAccept))
				assert.NotEmpty(t, req.Header.Get(headerAuth))
				assert.NotEmpty(t, req.URL.Query().Get("per_page"))
				assert.NotEmpty(t, req.URL.Query().Get("page"))
			}
		})
	}
}

func TestClient_NewUploadRequest(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		url           string
		filepath      string
		expectedError string
	}{
		{
			name:          "InvalidURL",
			ctx:           context.Background(),
			url:           ":invalid",
			filepath:      "",
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "NoFile",
			ctx:           context.Background(),
			url:           "/repos/octocat/Hello-World/releases/1/assets",
			filepath:      "",
			expectedError: `open : no such file or directory`,
		},
		{
			name:          "BadFile",
			ctx:           context.Background(),
			url:           "/repos/octocat/Hello-World/releases/1/assets",
			filepath:      "/dev/null",
			expectedError: `EOF`,
		},
		{
			name:          "NilContext",
			ctx:           nil,
			url:           "/repos/octocat/Hello-World/releases/1/assets",
			filepath:      "./fixture/asset",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "Success",
			ctx:           context.Background(),
			url:           "/repos/octocat/Hello-World/releases/1/assets",
			filepath:      "./fixture/asset",
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{
				uploadURL: publicUploadURL,
				authToken: "abcdefghijklmnopqrstuvwxyz",
			}

			req, closer, err := c.NewUploadRequest(tc.ctx, tc.url, tc.filepath)
			if err == nil {
				defer func() {
					assert.NoError(t, closer.Close())
				}()
			}

			if tc.expectedError != "" {
				assert.Nil(t, req)
				assert.Nil(t, closer)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.NotNil(t, closer)
				assert.NotEmpty(t, req.Header.Get(headerUserAgent))
				assert.NotEmpty(t, req.Header.Get(headerAccept))
				assert.NotEmpty(t, req.Header.Get(headerContentType))
				assert.NotEmpty(t, req.Header.Get(headerAuth))
			}
		})
	}
}

func TestClient_NewDownloadRequest(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		url           string
		expectedError string
	}{
		{
			name:          "InvalidURL",
			ctx:           context.Background(),
			url:           ":invalid",
			expectedError: `parse ":invalid": missing protocol scheme`,
		},
		{
			name:          "NilContext",
			ctx:           nil,
			url:           "/octocat/Hello-World/releases/download/v1.0.0/asset",
			expectedError: `net/http: nil Context`,
		},
		{
			name:          "Success",
			ctx:           context.Background(),
			url:           "/octocat/Hello-World/releases/download/v1.0.0/asset",
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{
				downloadURL: publicDownloadURL,
				authToken:   "abcdefghijklmnopqrstuvwxyz",
			}

			req, err := c.NewDownloadRequest(tc.ctx, tc.url)

			if tc.expectedError != "" {
				assert.Nil(t, req)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.NotEmpty(t, req.Header.Get(headerUserAgent))
				assert.NotEmpty(t, req.Header.Get(headerAuth))
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	type user struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	reset := time.Now().Add(time.Hour)

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		c                *Client
		reqMethod        string
		reqURL           string
		body             interface{}
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NoRemainingRateLimit",
			mockResponses: []MockResponse{},
			c: &Client{
				rates: map[rateGroup]Rate{
					rateGroupCore: {
						Limit:     5000,
						Used:      5000,
						Remaining: 0,
						Reset:     Epoch(reset.Unix()),
					},
				},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: rate limit 5000 used: rate limit will reset at ` + reset.Format("15:04:05"),
		},
		{
			name:          "HTTPClientError",
			mockResponses: []MockResponse{},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `Get "/user": unsupported protocol scheme ""`,
		},
		{
			name: "StatusBadRequest",
			mockResponses: []MockResponse{
				{"GET", "/user", 400, http.Header{}, `{
					"message": "Problems parsing JSON",
					"documentation_url": "https://docs.github.com/rest/reference/users#update-the-authenticated-user"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: 400 Problems parsing JSON`,
		},
		{
			name: "AuthError",
			mockResponses: []MockResponse{
				{"GET", "/user", 401, http.Header{}, `{
					"message": "Requires authentication",
					"documentation_url": "https://docs.github.com/rest/reference/users#get-the-authenticated-user"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: 401 Requires authentication`,
		},
		{
			name: "RateLimitError",
			mockResponses: []MockResponse{
				{
					"GET", "/user", 403,
					http.Header{
						headerRateRemaining: {"0"},
						headerRateReset:     {"1605125898"},
					},
					`{
						"message": "API rate limit exceeded",
						"documentation_url": "https://developer.github.com/v3/#rate-limiting"
					}`,
				},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: rate limit 0 used: rate limit will reset at ` + time.Unix(1605125898, 0).Format("15:04:05"),
		},
		{
			name: "RateLimitAbuseError",
			mockResponses: []MockResponse{
				{
					"GET", "/user", 403,
					http.Header{
						headerRetryAfter: {"30"},
					},
					`{
						"message": "You have triggered an abuse detection mechanism",
						"documentation_url": "https://developer.github.com/v3/#abuse-rate-limits"
					}`,
				},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: 403 You have triggered an abuse detection mechanism`,
		},
		{
			name: "NotFoundError",
			mockResponses: []MockResponse{
				{"GET", "/users/octocat", 404, http.Header{}, `{
					"message": "Not Found",
					"documentation_url": "https://docs.github.com/rest"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/users/octocat",
			body:          nil,
			expectedError: `GET /users/octocat: 404 Not Found`,
		},
		{
			name: "StatusInternalServerError",
			mockResponses: []MockResponse{
				{"GET", "/user", 500, http.Header{}, `Internal server error`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          nil,
			expectedError: `GET /user: 500 `,
		},
		{
			name: "InvalidJSON",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, http.Header{}, `{`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod:     "GET",
			reqURL:        "/user",
			body:          new(user),
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success_Writer",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, header, `{
						"login": "octocat",
						"id": 1,
						"name": "The Octocat",
						"email": "octocat@github.com"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod: "GET",
			reqURL:    "/user",
			body:      new(bytes.Buffer),
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
		{
			name: "Success_Struct",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, header, `{
						"login": "octocat",
						"id": 1,
						"name": "The Octocat",
						"email": "octocat@github.com"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod: "GET",
			reqURL:    "/user",
			body:      new(user),
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
		{
			name: "Success_Map",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, header, `{
						"login": "octocat",
						"id": 1,
						"name": "The Octocat",
						"email": "octocat@github.com"
				}`},
			},
			c: &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
			},
			reqMethod: "GET",
			reqURL:    "/user",
			body:      new(map[string]interface{}),
			expectedResponse: &Response{
				Pages: expectedPages,
				Rate:  expectedRate,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.mockResponses) > 0 {
				ts := newHTTPTestServer(tc.mockResponses...)
				defer ts.Close()

				serverURL, _ := url.Parse(ts.URL)

				reqURL, err := serverURL.Parse(tc.reqURL)
				assert.NoError(t, err)

				tc.reqURL = reqURL.String()
			}

			req, err := http.NewRequest(tc.reqMethod, tc.reqURL, nil)
			assert.NoError(t, err)

			// UAT
			resp, err := tc.c.Do(req, tc.body)

			if tc.expectedError != "" {
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestClient_EnsureScopes(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []MockResponse
		ctx           context.Context
		scopes        []Scope
		expectedError string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			ctx:           nil,
			scopes:        []Scope{ScopeRepo},
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 401, http.Header{}, `bad credentials`},
			},
			ctx:           context.Background(),
			scopes:        []Scope{ScopeRepo},
			expectedError: `HEAD /user: 401 `,
		},
		{
			name: "MissingScope",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, http.Header{}, ``},
			},
			ctx:           context.Background(),
			scopes:        []Scope{ScopeRepo},
			expectedError: `auth token does not have the scope: repo`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"HEAD", "/user", 200, http.Header{
					"X-OAuth-Scopes": []string{"repo"},
				}, ``},
			},
			ctx:           context.Background(),
			scopes:        []Scope{ScopeRepo},
			expectedError: ``,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{
				httpClient: &http.Client{},
				rates:      map[rateGroup]Rate{},
				apiURL:     publicUploadURL,
			}

			ts := newHTTPTestServer(tc.mockResponses...)
			defer ts.Close()

			c.apiURL, _ = url.Parse(ts.URL)

			err := c.EnsureScopes(tc.ctx, tc.scopes...)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestClient_Repo(t *testing.T) {
	tests := []struct {
		name          string
		owner         string
		repo          string
		expectedError string
	}{
		{
			name:  "OK",
			owner: "octocat",
			repo:  "Hello-World",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{}

			repo := c.Repo(tc.owner, tc.repo)

			assert.NotNil(t, repo)
			assert.Equal(t, c, repo.client)
			assert.Equal(t, tc.owner, repo.owner)
			assert.Equal(t, tc.repo, repo.repo)

			assert.NotNil(t, repo.Pulls)
			assert.Equal(t, c, repo.Pulls.client)
			assert.Equal(t, tc.owner, repo.Pulls.owner)
			assert.Equal(t, tc.repo, repo.Pulls.repo)

			assert.NotNil(t, repo.Issues)
			assert.Equal(t, c, repo.Issues.client)
			assert.Equal(t, tc.owner, repo.Issues.owner)
			assert.Equal(t, tc.repo, repo.Issues.repo)

			assert.NotNil(t, repo.Releases)
			assert.Equal(t, c, repo.Releases.client)
			assert.Equal(t, tc.owner, repo.Releases.owner)
			assert.Equal(t, tc.repo, repo.Releases.repo)
		})
	}
}
