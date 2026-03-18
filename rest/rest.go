// package rest provides types and methods for calling the GitHub REST API.
package rest

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	relFirstRE = regexp.MustCompile(`<https://api.github.com/[^>]+[?&]page=(\d+)[^>]*>; rel="first"`)
	relPrevRE  = regexp.MustCompile(`<https://api.github.com/[^>]+[?&]page=(\d+)[^>]*>; rel="prev"`)
	relNextRE  = regexp.MustCompile(`<https://api.github.com/[^>]+[?&]page=(\d+)[^>]*>; rel="next"`)
	relLastRE  = regexp.MustCompile(`<https://api.github.com/[^>]+[?&]page=(\d+)[^>]*>; rel="last"`)
)

const (
	headerLink          = "Link"
	headerRateResource  = "X-Ratelimit-Resource"
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

	r.Rate.Resource = h.Get(headerRateResource)

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
type Pages struct {
	First int
	Prev  int
	Next  int
	Last  int
}

// Rate represents the rate limit status for the authenticated user.
type Rate struct {
	// The resource being rate limited.
	Resource string `json:"resource"`
	// The number of requests per hour.
	Limit int `json:"limit"`
	// The number of requests used in the current hour.
	Used int `json:"used"`
	// The number of requests remaining in the current hour.
	Remaining int `json:"remaining"`
	// The time at which the current rate will reset.
	Reset Epoch `json:"reset"`
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

// rateGroup determines the rate limit group for GitHub REST API.
type rateGroup string

const (
	rateGroupCore    = rateGroup("core")
	rateGroupSearch  = rateGroup("search")
	rateGroupGraphQL = rateGroup("graphql")
)

func getRateGroup(u *url.URL) rateGroup {
	switch {
	case strings.HasPrefix(u.Path, "/search"):
		return rateGroupSearch
	case strings.HasPrefix(u.Path, "/graphql"):
		return rateGroupGraphQL
	default:
		return rateGroupCore
	}
}
