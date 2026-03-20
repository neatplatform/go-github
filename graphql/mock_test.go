package graphql

import (
	"context"
	"net/http"

	"github.com/neatplatform/go-github"
)

type (
	MockGithubClient struct {
		NewRequestIndex int
		NewRequestMocks []NewRequestMock

		DoIndex int
		DoMocks []DoMock
	}

	NewRequestMock struct {
		Func       func(context.Context, string, string, any) (*http.Request, error)
		InCtx      context.Context
		InMethod   string
		InURL      string
		InBody     any
		OutRequest *http.Request
		OutError   error
	}

	DoMock struct {
		Func        func(*http.Request, any) (*github.Response, error)
		InReq       any
		InBody      any
		OutResponse *github.Response
		OutError    error
	}
)

func (m *MockGithubClient) NewRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	if m.NewRequestIndex >= len(m.NewRequestMocks) {
		panic("NewRequest called more times than expected")
	}

	i := m.NewRequestIndex
	m.NewRequestIndex++

	// Use the custom function if provided.
	if m.NewRequestMocks[i].Func != nil {
		return m.NewRequestMocks[i].Func(ctx, method, url, body)
	}

	// Otherwise, record the inputs and return the predefined outputs.
	m.NewRequestMocks[i].InCtx = ctx
	m.NewRequestMocks[i].InMethod = method
	m.NewRequestMocks[i].InURL = url
	m.NewRequestMocks[i].InBody = body
	return m.NewRequestMocks[i].OutRequest, m.NewRequestMocks[i].OutError
}

func (m *MockGithubClient) Do(req *http.Request, body any) (*github.Response, error) {
	if m.DoIndex >= len(m.DoMocks) {
		panic("Do called more times than expected")
	}

	i := m.DoIndex
	m.DoIndex++

	// Use the custom function if provided.
	if m.DoMocks[i].Func != nil {
		return m.DoMocks[i].Func(req, body)
	}

	// Otherwise, record the inputs and return the predefined outputs.
	m.DoMocks[i].InReq = req
	m.DoMocks[i].InBody = body
	return m.DoMocks[i].OutResponse, m.DoMocks[i].OutError
}
