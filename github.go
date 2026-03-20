// Package github provides types and methods for calling the GitHub REST API.
package github

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// The Link-header pagination regexes should handle both the public and enterprise GitHub instances.
	relFirstRE = regexp.MustCompile(`<https://[^>]+/[^>]+[?&]page=(\d+)[^>]*>; rel="first"`)
	relPrevRE  = regexp.MustCompile(`<https://[^>]+/[^>]+[?&]page=(\d+)[^>]*>; rel="prev"`)
	relNextRE  = regexp.MustCompile(`<https://[^>]+/[^>]+[?&]page=(\d+)[^>]*>; rel="next"`)
	relLastRE  = regexp.MustCompile(`<https://[^>]+/[^>]+[?&]page=(\d+)[^>]*>; rel="last"`)
)

const (
	headerLink          = "Link"
	headerRateResource  = "X-RateLimit-Resource"
	headerRateLimit     = "X-RateLimit-Limit"
	headerRateUsed      = "X-RateLimit-Used"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

// Response represents an HTTP response for GitHub REST API.
type Response struct {
	*http.Response

	Pages Pages
	Rate  Rate
}

func newResponse(resp *http.Response) *Response {
	r := &Response{
		Response: resp,
	}

	h := resp.Header

	if link := h.Get(headerLink); link != "" {
		if m := relFirstRE.FindStringSubmatch(link); len(m) == 2 {
			r.Pages.First, _ = strconv.Atoi(m[1])
		}

		if m := relPrevRE.FindStringSubmatch(link); len(m) == 2 {
			r.Pages.Prev, _ = strconv.Atoi(m[1])
		}

		if m := relNextRE.FindStringSubmatch(link); len(m) == 2 {
			r.Pages.Next, _ = strconv.Atoi(m[1])
		}

		if m := relLastRE.FindStringSubmatch(link); len(m) == 2 {
			r.Pages.Last, _ = strconv.Atoi(m[1])
		}
	}

	r.Rate.Resource = Category(h.Get(headerRateResource))

	if limit := h.Get(headerRateLimit); limit != "" {
		r.Rate.Limit, _ = strconv.Atoi(limit)
	}

	if used := h.Get(headerRateUsed); used != "" {
		r.Rate.Used, _ = strconv.Atoi(used)
	}

	if remaining := h.Get(headerRateRemaining); remaining != "" {
		r.Rate.Remaining, _ = strconv.Atoi(remaining)
	}

	if reset := h.Get(headerRateReset); reset != "" {
		i64, _ := strconv.ParseInt(reset, 10, 64)
		r.Rate.Reset = Epoch(i64)
	}

	return r
}

// Pages represents the pagination information for GitHub REST API.
//
// See https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api
type Pages struct {
	First int
	Prev  int
	Next  int
	Last  int
}

// Rate represents the rate limit status for the authenticated user.
//
// See https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
type Rate struct {
	// The resource being rate limited.
	Resource Category `json:"resource"`
	// The number of requests per hour.
	Limit int `json:"limit"`
	// The number of requests used in the current hour.
	Used int `json:"used"`
	// The number of requests remaining in the current hour.
	Remaining int `json:"remaining"`
	// The time at which the current rate will reset.
	Reset Epoch `json:"reset"`
}

// Category determines the rate limit resource category for GitHub REST API.
//
// See https://docs.github.com/en/rest/rate-limit/rate-limit#get-rate-limit-status-for-the-authenticated-user
type Category string

const (
	CategoryCore       = Category("core")
	CategorySearch     = Category("search")
	CategoryCodeSearch = Category("code_search")
	CategoryGraphQL    = Category("graphql")
)

func getCategory(path string) Category {
	switch {
	case strings.HasPrefix(path, "/search/code"):
		return CategoryCodeSearch
	case strings.HasPrefix(path, "/search"):
		return CategorySearch
	case strings.HasPrefix(path, "/graphql"):
		return CategoryGraphQL
	default:
		return CategoryCore
	}
}

// Epoch is a Unix timestamp.
type Epoch int64

// Time returns the Time representation of an epoch timestamp.
func (e Epoch) Time() time.Time {
	return time.Unix(int64(e), 0)
}

// String returns string representation of an epoch timestamp.
func (e Epoch) String() string {
	return e.Time().Format("15:04:05")
}
