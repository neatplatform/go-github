package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	headerAuth        = "Authorization"
	headerUserAgent   = "User-Agent"
	headerContentType = "Content-Type"
	headerAccept      = "Accept"
	headerScopes      = "X-OAuth-Scopes"
	headerRetryAfter  = "Retry-After"
)

const (
	mediaJSON   = "application/json"
	mediaTypeV3 = "application/vnd.github.v3+json"
)

var (
	userAgent = filepath.Base(os.Args[0])

	publicAPIURL, _      = url.Parse("https://api.github.com")
	publicUploadURL, _   = url.Parse("https://uploads.github.com")
	publicDownloadURL, _ = url.Parse("https://github.com")
)

// Client is used for making API calls to GitHub REST API.
type Client struct {
	httpClient *http.Client
	ratesMutex sync.Mutex
	rates      map[rateGroup]Rate

	apiURL      *url.URL
	uploadURL   *url.URL
	downloadURL *url.URL
	authToken   string

	// Services
	Users  *UserService
	Search *SearchService
}

func newHTTPClient() *http.Client {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	return client
}

// NewClient creates a new client for calling public GitHub REST API.
func NewClient(authToken string) *Client {
	c := &Client{
		httpClient:  newHTTPClient(),
		rates:       map[rateGroup]Rate{},
		apiURL:      publicAPIURL,
		uploadURL:   publicUploadURL,
		downloadURL: publicDownloadURL,
		authToken:   authToken,
	}

	c.Users = &UserService{
		client: c,
	}

	c.Search = &SearchService{
		client: c,
	}

	return c
}

// NewEnterpriseClient creates a new client for calling an enterprise GitHub REST API.
func NewEnterpriseClient(apiURL, uploadURL, downloadURL, authToken string) (*Client, error) {
	entAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	entUploadURL, err := url.Parse(uploadURL)
	if err != nil {
		return nil, err
	}

	entDownloadURL, err := url.Parse(downloadURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient:  newHTTPClient(),
		rates:       map[rateGroup]Rate{},
		apiURL:      entAPIURL,
		uploadURL:   entUploadURL,
		downloadURL: entDownloadURL,
		authToken:   authToken,
	}

	c.Users = &UserService{
		client: c,
	}

	c.Search = &SearchService{
		client: c,
	}

	return c, nil
}

// NewRequest creates a new HTTP request for a GitHub REST API.
// If body implements the io.Reader interface, the raw request body will be read.
// Otherwise, the request body will be JOSN-encoded.
func (c *Client) NewRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	u, err := c.apiURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			reader = r
		} else {
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(body); err != nil {
				return nil, err
			}
			reader = buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerUserAgent, userAgent)
	req.Header.Set(headerAccept, mediaTypeV3)

	if c.authToken != "" {
		req.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", c.authToken))
	}

	if body != nil {
		req.Header.Set(headerContentType, mediaJSON)
	}

	return req, nil
}

// NewPageRequest creates a new HTTP request for a GitHub REST API with page parameters.
// If body implements the io.Reader interface, the raw request body will be read.
// Otherwise, the request body will be JOSN-encoded.
func (c *Client) NewPageRequest(ctx context.Context, method, url string, pageSize, pageNo int, body interface{}) (*http.Request, error) {
	req, err := c.NewRequest(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if pageSize > 0 {
		q.Add("per_page", strconv.Itoa(pageSize))
	}
	if pageNo > 0 {
		q.Add("page", strconv.Itoa(pageNo))
	}

	req.URL.RawQuery = q.Encode()

	return req, nil
}

// NewUploadRequest creates a new HTTP request for uploading a file to a GitHub release.
// When successful, it returns a closer for the given file that should be closed after making the request.
func (c *Client) NewUploadRequest(ctx context.Context, url, filepath string) (*http.Request, io.Closer, error) {
	u, err := c.uploadURL.Parse(url)
	if err != nil {
		return nil, nil, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	// Read the first 512 bytes of file to determine the media type of the file
	buff := make([]byte, 512)
	if _, err := f.Read(buff); err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	// http.DetectContentType will return "application/octet-stream" if it cannot determine a more specific one
	mediaType := http.DetectContentType(buff)

	// Reset the offset back to the beginning of the file
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), f)
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	req.ContentLength = stat.Size()
	req.Header.Set(headerUserAgent, userAgent)
	req.Header.Set(headerAccept, mediaTypeV3)
	req.Header.Set(headerContentType, mediaType)

	if c.authToken != "" {
		req.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", c.authToken))
	}

	return req, f, nil
}

// NewDownloadRequest creates a new HTTP request for downloading a file from a GitHub release.
func (c *Client) NewDownloadRequest(ctx context.Context, url string) (*http.Request, error) {
	u, err := c.downloadURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerUserAgent, userAgent)

	if c.authToken != "" {
		req.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", c.authToken))
	}

	return req, nil
}

// Do makes an HTTP request and returns the API response.
// If body implements the io.Writer interface, the raw response body will be copied to.
// Otherwise, the response body will be JOSN-decoded into it.
func (c *Client) Do(req *http.Request, body interface{}) (*Response, error) {
	/* -------------------- CHECK RATE LIMITS -------------------- */

	g := getRateGroup(req.URL)

	c.ratesMutex.Lock()
	rate, ok := c.rates[g]
	c.ratesMutex.Unlock()

	if ok && rate.Remaining == 0 && time.Now().Before(rate.Reset.Time()) {
		return nil, &RateLimitError{
			Request: req,
			Rate:    rate,
		}
	}

	/* -------------------- MAKE THE REQUEST -------------------- */

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Ensure we fully read and close the response body, so the underlying TCP connection can be reused.
		// If it errors, the TCP connection will not be reused anyway.
		_, _ = io.Copy(io.Discard, r.Body)
		_ = r.Body.Close()
	}()

	resp := newResponse(r)

	// Update rate limits
	c.ratesMutex.Lock()
	c.rates[g] = resp.Rate
	c.ratesMutex.Unlock()

	/* -------------------- CHECK THE RESPONSE -------------------- */

	isSuccess := func(statusCode int) bool {
		return statusCode == http.StatusOK ||
			statusCode == http.StatusCreated ||
			statusCode == http.StatusNoContent
	}

	if !isSuccess(r.StatusCode) {
		respErr := &ResponseError{
			Response: r,
		}

		b, err := io.ReadAll(resp.Body)
		if err == nil && b != nil {
			_ = json.Unmarshal(b, respErr)
		}

		// Restore response body
		// r.Body = io.NopCloser(bytes.NewBuffer(b))

		switch r.StatusCode {
		case http.StatusBadRequest:
			return nil, respErr

		case http.StatusUnauthorized:
			return nil, &AuthError{
				err: respErr,
			}

		case http.StatusForbidden:
			if r.Header.Get(headerRateRemaining) == "0" {
				return nil, &RateLimitError{
					err:     respErr,
					Request: req,
					Rate:    resp.Rate,
				}
			} else if strings.HasSuffix(respErr.DocumentationURL, "#abuse-rate-limits") {
				retryAfter, _ := time.ParseDuration(r.Header.Get(headerRetryAfter) + "s")
				return nil, &RateLimitAbuseError{
					err:        respErr,
					Rate:       resp.Rate,
					RetryAfter: retryAfter,
				}
			}

		case http.StatusNotFound:
			return nil, &NotFoundError{
				err: respErr,
			}

		default:
			return nil, respErr
		}
	}

	/* -------------------- READ THE BODY -------------------- */

	if body != nil {
		if w, ok := body.(io.Writer); ok {
			if _, err := io.Copy(w, r.Body); err != nil {
				return nil, err
			}
		} else {
			if err := json.NewDecoder(r.Body).Decode(body); err != nil && err != io.EOF {
				return nil, err
			}
		}
	}

	return resp, nil
}

// EnsureScopes makes sure the client and the auth token have the given scopes.
//
// See https://docs.github.com/developers/apps/scopes-for-oauth-apps
func (c *Client) EnsureScopes(ctx context.Context, scopes ...Scope) error {
	// Call an endpoint to get the OAuth scopes of the auth token from the headers
	req, err := c.NewRequest(ctx, "HEAD", "/user", nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		return err
	}

	// Ensure the auth token has all the required OAuth scopes
	oauthScopes := resp.Header.Get(headerScopes)
	for _, scope := range scopes {
		if !strings.Contains(oauthScopes, string(scope)) {
			return fmt.Errorf("auth token does not have the scope: %s", scope)
		}
	}

	return nil
}

// Repo returns a service providing GitHub REST APIs for a specific repository.
func (c *Client) Repo(owner, repo string) *RepoService {
	return &RepoService{
		client: c,
		owner:  owner,
		repo:   repo,
		Pulls: &PullService{
			client: c,
			owner:  owner,
			repo:   repo,
		},
		Issues: &IssueService{
			client: c,
			owner:  owner,
			repo:   repo,
		},
		Releases: &ReleaseService{
			client: c,
			owner:  owner,
			repo:   repo,
		},
	}
}
