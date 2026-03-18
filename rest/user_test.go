package rest

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	userBody = `{
		"login": "octocat",
		"id": 1,
		"url": "https://api.github.com/users/octocat",
		"html_url": "https://github.com/octocat",
		"type": "User",
		"site_admin": false,
		"name": "The Octocat",
		"email": "octocat@github.com"
	}`
)

var (
	user = User{
		ID:      1,
		Login:   "octocat",
		Type:    "User",
		Email:   "octocat@github.com",
		Name:    "The Octocat",
		URL:     "https://api.github.com/users/octocat",
		HTMLURL: "https://github.com/octocat",
	}
)

func TestUserService_User(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicUploadURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *UserService
		ctx              context.Context
		expectedUser     *User
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &UserService{
				client: c,
			},
			ctx:           nil,
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/user", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &UserService{
				client: c,
			},
			ctx:           context.Background(),
			expectedError: `GET /user: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, http.Header{}, `{`},
			},
			s: &UserService{
				client: c,
			},
			ctx:           context.Background(),
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/user", 200, header, userBody},
			},
			s: &UserService{
				client: c,
			},
			ctx:          context.Background(),
			expectedUser: &user,
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

			user, resp, err := tc.s.User(tc.ctx)

			if tc.expectedError != "" {
				assert.Nil(t, user)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{},
		rates:      map[rateGroup]Rate{},
		apiURL:     publicUploadURL,
	}

	tests := []struct {
		name             string
		mockResponses    []MockResponse
		s                *UserService
		ctx              context.Context
		username         string
		expectedUser     *User
		expectedResponse *Response
		expectedError    string
	}{
		{
			name:          "NilContext",
			mockResponses: []MockResponse{},
			s: &UserService{
				client: c,
			},
			ctx:           nil,
			username:      "octocat",
			expectedError: `net/http: nil Context`,
		},
		{
			name: "InvalidStatusCode",
			mockResponses: []MockResponse{
				{"GET", "/users/octocat", 401, http.Header{}, `{
					"message": "Bad credentials"
				}`},
			},
			s: &UserService{
				client: c,
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: `GET /users/octocat: 401 Bad credentials`,
		},
		{
			name: "ّInvalidResponse",
			mockResponses: []MockResponse{
				{"GET", "/users/octocat", 200, http.Header{}, `{`},
			},
			s: &UserService{
				client: c,
			},
			ctx:           context.Background(),
			username:      "octocat",
			expectedError: `unexpected EOF`,
		},
		{
			name: "Success",
			mockResponses: []MockResponse{
				{"GET", "/users/octocat", 200, header, userBody},
			},
			s: &UserService{
				client: c,
			},
			ctx:          context.Background(),
			username:     "octocat",
			expectedUser: &user,
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

			user, resp, err := tc.s.Get(tc.ctx, tc.username)

			if tc.expectedError != "" {
				assert.Nil(t, user)
				assert.Nil(t, resp)
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Response)
				assert.Equal(t, tc.expectedResponse.Pages, resp.Pages)
				assert.Equal(t, tc.expectedResponse.Rate, resp.Rate)
			}
		})
	}
}
