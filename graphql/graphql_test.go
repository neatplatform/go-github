package graphql

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/neatplatform/go-github"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		c    *github.Client
	}{
		{
			name: "OK",
			c:    github.NewClient("abcdefghijklmnopqrstuvwxyz"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := New(tc.c)

			assert.NotNil(t, g)
			assert.Equal(t, tc.c, g.client)
		})
	}
}

func TestGraphQL_Query(t *testing.T) {
	tests := []struct {
		name           string
		g              *GraphQL
		ctx            context.Context
		query          string
		vars           map[string]any
		result         any
		expectedError  string
		expectedResult any
	}{
		{
			name: "NewRequestError",
			g: &GraphQL{
				client: &MockGithubClient{
					NewRequestMocks: []NewRequestMock{
						{OutError: errors.New("dummy error")},
					},
				},
			},
			ctx: nil,
			query: `query {
				viewer {
					login
				}
			}`,
			vars:          nil,
			result:        nil,
			expectedError: `failed to create GraphQL request: dummy error`,
		},
		{
			name: "DoError",
			g: &GraphQL{
				client: &MockGithubClient{
					NewRequestMocks: []NewRequestMock{
						{OutRequest: &http.Request{}},
					},
					DoMocks: []DoMock{
						{OutError: errors.New("dummy error")},
					},
				},
			},
			ctx: nil,
			query: `query {
				viewer {
					login
				}
			}`,
			vars:          nil,
			result:        nil,
			expectedError: `failed to execute GraphQL request: dummy error`,
		},
		{
			name: "GraphQLErrors",
			g: &GraphQL{
				client: &MockGithubClient{
					NewRequestMocks: []NewRequestMock{
						{OutRequest: &http.Request{}},
					},
					DoMocks: []DoMock{
						{
							Func: func(req *http.Request, body any) (*github.Response, error) {
								if v, ok := body.(*Response); ok {
									v.Errors = []Error{
										{Message: "A query attribute must be specified and must be a string."},
									}
								}

								return &github.Response{
									Response: &http.Response{},
								}, nil
							},
						},
					},
				},
			},
			ctx: nil,
			query: `query {
				viewer {
					login
				}
			}`,
			vars:          nil,
			result:        nil,
			expectedError: "GraphQL Error:\nA query attribute must be specified and must be a string.",
		},
		{
			name: "ResultNotPointer",
			g: &GraphQL{
				client: &MockGithubClient{
					NewRequestMocks: []NewRequestMock{
						{OutRequest: &http.Request{}},
					},
					DoMocks: []DoMock{
						{
							Func: func(req *http.Request, body any) (*github.Response, error) {
								if v, ok := body.(*Response); ok {
									v.Data = []byte(`
										{
											"viewer": {
												"login": "octocat"
											}
										}
									`)
								}

								return &github.Response{
									Response: &http.Response{},
								}, nil
							},
						},
					},
				},
			},
			ctx: nil,
			query: `query {
				viewer {
					login
				}
			}`,
			vars:          nil,
			result:        map[string]any{},
			expectedError: `json: Unmarshal(non-pointer map[string]interface {})`,
		},
		{
			name: "Success",
			g: &GraphQL{
				client: &MockGithubClient{
					NewRequestMocks: []NewRequestMock{
						{OutRequest: &http.Request{}},
					},
					DoMocks: []DoMock{
						{
							Func: func(req *http.Request, body any) (*github.Response, error) {
								if v, ok := body.(*Response); ok {
									v.Data = []byte(`
										{
											"viewer": {
												"login": "octocat"
											}
										}
									`)
								}

								return &github.Response{
									Response: &http.Response{},
								}, nil
							},
						},
					},
				},
			},
			ctx: nil,
			query: `query {
				viewer {
					login
				}
			}`,
			vars:   nil,
			result: &map[string]any{},
			expectedResult: &map[string]any{
				"viewer": map[string]any{
					"login": "octocat",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.g.Query(tc.ctx, tc.query, tc.vars, tc.result)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, tc.result)
			}
		})
	}
}
