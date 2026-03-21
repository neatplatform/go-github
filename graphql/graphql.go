// Package graphql provides types and methods for calling the GitHub GraphQL API.
package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/neatplatform/go-github"
)

type (
	// Request represents a GraphQL request payload.
	Request struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables,omitempty"`
	}

	// Response represents a GraphQL response payload.
	Response struct {
		Data   json.RawMessage `json:"data"`
		Errors []Error         `json:"errors"`
	}

	// Error represents an error returned by the GraphQL API.
	Error struct {
		Message string `json:"message"`
	}

	// RateLimit represents the GraphQL API rate limit information.
	RateLimit struct {
		Limit     int       `json:"limit"`
		Remaining int       `json:"remaining"`
		ResetAt   time.Time `json:"resetAt"` // RFC3339
		Cost      int       `json:"cost"`
	}

	// PageInfo represents the pagination information for GraphQL queries.
	PageInfo struct {
		StartCursor     string `json:"startCursor"`
		EndCursor       string `json:"endCursor"`
		HasNextPage     bool   `json:"hasNextPage"`
		HasPreviousPage bool   `json:"hasPreviousPage"`
	}
)

// GraphQL is used for making calls to GitHub GraphQL API.
type GraphQL struct {
	client githubClient
}

// Define an interface for dependency injection and mocking in tests.
type githubClient interface {
	NewRequest(context.Context, string, string, any) (*http.Request, error)
	Do(*http.Request, any) (*github.Response, error)
}

// New creates a new GraphQL client with the provided GitHub REST client.
func New(c *github.Client) *GraphQL {
	return &GraphQL{
		client: c,
	}
}

// Query executes a GraphQL query.
func (g *GraphQL) Query(ctx context.Context, query string, vars map[string]any, result any) error {
	req, err := g.client.NewRequest(ctx, "POST", "/graphql", Request{
		Query:     query,
		Variables: vars,
	})

	if err != nil {
		return fmt.Errorf("failed to create GraphQL request: %w", err)
	}

	var body Response
	if _, err := g.client.Do(req, &body); err != nil {
		return fmt.Errorf("failed to execute GraphQL request: %w", err)
	}

	var errs error
	for _, e := range body.Errors {
		errs = errors.Join(errs, errors.New(e.Message))
	}

	if errs != nil {
		return fmt.Errorf("GraphQL Error:\n%w", errs)
	}

	if result != nil {
		if err := json.Unmarshal(body.Data, result); err != nil {
			return err
		}
	}

	return nil
}
