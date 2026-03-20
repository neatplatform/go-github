package github

import (
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func parseGitHubTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return t
}

func parseGitHubTimePtr(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return &t
}

type MockResponse struct {
	Method             string
	Path               string
	ResponseStatusCode int
	ResponseHeader     http.Header
	ResponseBody       string
}

func newHTTPTestServer(mocks ...MockResponse) *httptest.Server {
	mux := http.NewServeMux()

	for _, m := range mocks {
		m := m

		mux.HandleFunc(m.Path, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m.Method {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			for k, vals := range m.ResponseHeader {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}

			w.WriteHeader(m.ResponseStatusCode)
			_, _ = io.WriteString(w, m.ResponseBody)
		})
	}

	return httptest.NewServer(mux)
}
