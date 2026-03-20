package github

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseErrorError(t *testing.T) {
	req, _ := http.NewRequest("PATCH", "/user", nil)

	tests := []struct {
		name          string
		err           *ResponseError
		expectedError string
	}{
		{
			name: "OK",
			err: &ResponseError{
				Response: &http.Response{
					StatusCode: 400,
					Request:    req,
				},
				Status:           "400",
				Message:          "Problems parsing JSON",
				DocumentationURL: "https://docs.github.com/rest/reference/users#update-the-authenticated-user",
			},
			expectedError: "PATCH /user: [400] Problems parsing JSON",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualError(t, tc.err, tc.expectedError)
		})
	}
}

func TestPrimaryRateLimitError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)

	tests := []struct {
		name          string
		err           *PrimaryRateLimitError
		expectedError string
	}{
		{
			name: "WithoutResponseError",
			err: &PrimaryRateLimitError{
				Rate: Rate{
					Limit:     5000,
					Used:      5000,
					Remaining: 0,
					Reset:     Epoch(1605125898),
				},
			},
			expectedError: "primary rate limit exhausted, resets at " + time.Unix(1605125898, 0).Format("15:04:05"),
		},
		{
			name: "WithResponseError",
			err: &PrimaryRateLimitError{
				Err: &ResponseError{
					Response: &http.Response{
						StatusCode: 403,
						Request:    req,
					},
					Status:           "403",
					Message:          "API rate limit exceeded",
					DocumentationURL: "https://developer.github.com/v3/#rate-limiting",
				},
				Rate: Rate{
					Limit:     5000,
					Used:      5000,
					Remaining: 0,
					Reset:     Epoch(1605125898),
				},
			},
			expectedError: "GET /user: [403] API rate limit exceeded: primary rate limit exhausted, resets at " + time.Unix(1605125898, 0).Format("15:04:05"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualError(t, tc.err, tc.expectedError)
			assert.Equal(t, tc.err.Err, tc.err.Unwrap())
		})
	}
}

func TestSecondaryRateLimitError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)

	tests := []struct {
		name          string
		err           *SecondaryRateLimitError
		expectedError string
	}{
		{
			name: "WithoutResponseError",
			err: &SecondaryRateLimitError{
				Rate: Rate{
					Limit:     5000,
					Used:      1000,
					Remaining: 4000,
					Reset:     Epoch(1605125898),
				},
				RetryAfter: 30 * time.Second,
			},
			expectedError: "secondary rate limit exceeded, retry after 30s",
		},
		{
			name: "WithResponseError",
			err: &SecondaryRateLimitError{
				Err: &ResponseError{
					Response: &http.Response{
						StatusCode: 403,
						Request:    req,
					},
					Status:           "403",
					Message:          "You have triggered an abuse detection mechanism",
					DocumentationURL: "https://developer.github.com/v3/#abuse-rate-limits",
				},
				Rate: Rate{
					Limit:     5000,
					Used:      1000,
					Remaining: 4000,
					Reset:     Epoch(1605125898),
				},
				RetryAfter: 30 * time.Second,
			},
			expectedError: "GET /user: [403] You have triggered an abuse detection mechanism: secondary rate limit exceeded, retry after 30s",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualError(t, tc.err, tc.expectedError)
			assert.Equal(t, tc.err.Err, tc.err.Unwrap())
		})
	}
}
