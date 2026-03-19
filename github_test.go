package github

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	tests := []struct {
		name             string
		respHeader       http.Header
		expectedResponse *Response
	}{
		{
			name: "WithNextAndLast",
			respHeader: http.Header{
				headerLink:          {`<https://api.github.com/repos/octocat/Hello-World/issues?page=2&state=closed>; rel="next", <https://api.github.com/repos/octocat/Hello-World/issues?page=6&state=closed>; rel="last"`},
				headerRateResource:  {"core"},
				headerRateLimit:     {"5000"},
				headerRateUsed:      {"5"},
				headerRateRemaining: {"4995"},
				headerRateReset:     {"1605083281"},
			},
			expectedResponse: &Response{
				Pages: Pages{
					Next: 2,
					Last: 6,
				},
				Rate: Rate{
					Resource:  "core",
					Limit:     5000,
					Used:      5,
					Remaining: 4995,
					Reset:     Epoch(1605083281),
				},
			},
		},
		{
			name: "WithAll",
			respHeader: http.Header{
				headerLink:          {`<https://api.github.com/repos/octocat/Hello-World/issues?page=2&state=closed>; rel="prev", <https://api.github.com/repos/octocat/Hello-World/issues?page=4&state=closed>; rel="next", <https://api.github.com/repos/octocat/Hello-World/issues?page=6&state=closed>; rel="last", <https://api.github.com/repos/octocat/Hello-World/issues?page=1&state=closed>; rel="first"`},
				headerRateResource:  {"core"},
				headerRateLimit:     {"5000"},
				headerRateUsed:      {"10"},
				headerRateRemaining: {"4990"},
				headerRateReset:     {"1605083281"},
			},
			expectedResponse: &Response{
				Pages: Pages{
					First: 1,
					Prev:  2,
					Next:  4,
					Last:  6,
				},
				Rate: Rate{
					Resource:  "core",
					Limit:     5000,
					Used:      10,
					Remaining: 4990,
					Reset:     Epoch(1605083281),
				},
			},
		},
		{
			name: "WithPrevAndFirst",
			respHeader: http.Header{
				headerLink:          {`<https://api.github.com/repos/octocat/Hello-World/issues?page=5&state=closed>; rel="prev", <https://api.github.com/repos/octocat/Hello-World/issues?page=1&state=closed>; rel="first"`},
				headerRateResource:  {"core"},
				headerRateLimit:     {"5000"},
				headerRateUsed:      {"100"},
				headerRateRemaining: {"4900"},
				headerRateReset:     {"1605083281"},
			},
			expectedResponse: &Response{
				Pages: Pages{
					First: 1,
					Prev:  5,
				},
				Rate: Rate{
					Resource:  "core",
					Limit:     5000,
					Used:      100,
					Remaining: 4900,
					Reset:     Epoch(1605083281),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Header: http.Header{},
			}

			for k, vals := range tc.respHeader {
				for _, v := range vals {
					resp.Header.Add(k, v)
				}
			}

			r := newResponse(resp)

			assert.NotNil(t, r)
			tc.expectedResponse.Response = resp
			assert.Equal(t, tc.expectedResponse, r)
		})
	}
}

func TestEpoch(t *testing.T) {
	tests := []struct {
		name           string
		e              Epoch
		expectedTime   time.Time
		expectedString string
	}{
		{
			name: "OK",
			e:    Epoch(1605064490),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotZero(t, tc.e.Time())
			assert.NotEmpty(t, tc.e.String())
		})
	}
}

func TestGetCategory(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		expectedCategory Category
	}{
		{
			name:             "Core",
			path:             "/users/octocat",
			expectedCategory: CategoryCore,
		},
		{
			name:             "Search",
			path:             "/search",
			expectedCategory: CategorySearch,
		},
		{
			name:             "CodeSearch",
			path:             "/search/code",
			expectedCategory: CategoryCodeSearch,
		},
		{
			name:             "GraphQL",
			path:             "/graphql",
			expectedCategory: CategoryGraphQL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCategory, getCategory(tc.path))
		})
	}
}
