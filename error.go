package github

import (
	"fmt"
	"net/http"
	"time"
)

// ResponseError represents an error response from GitHub REST API.
type ResponseError struct {
	Response         *http.Response `json:"-"`
	Status           string         `json:"status"`
	Message          string         `json:"message"`
	DocumentationURL string         `json:"documentation_url,omitempty"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s %s: [%s] %s", e.Response.Request.Method, e.Response.Request.URL.Path, e.Status, e.Message)
}

// PrimaryRateLimitError is returned when the authenticated user has exhausted their primary API rate limit.
//
// See https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api#about-primary-rate-limits
type PrimaryRateLimitError struct {
	Err  *ResponseError
	Rate Rate
}

func (e *PrimaryRateLimitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("primary rate limit exhausted, resets at %s", e.Rate.Reset)
	}

	return fmt.Sprintf("%s: primary rate limit exhausted, resets at %s", e.Err.Error(), e.Rate.Reset)
}

func (e *PrimaryRateLimitError) Unwrap() error {
	return e.Err
}

// SecondaryRateLimitError is returned when the API detects abusive behavior or excessive requests that violate rate limit best practices.
//
// See https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api#about-secondary-rate-limits
type SecondaryRateLimitError struct {
	Err        *ResponseError
	Rate       Rate
	RetryAfter time.Duration
}

func (e *SecondaryRateLimitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("secondary rate limit exceeded, retry after %s", e.RetryAfter)
	}

	return fmt.Sprintf("%s: secondary rate limit exceeded, retry after %s", e.Err.Error(), e.RetryAfter)
}

func (e *SecondaryRateLimitError) Unwrap() error {
	return e.Err
}
