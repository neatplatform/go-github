package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neatplatform/go-github"
	"github.com/neatplatform/go-github/graphql"
)

func main() {
	const query = `query($owner: String!, $repo: String!, $cursor: String) {
		rateLimit {
			limit
			remaining
			resetAt
			cost
		}
		repository(owner: $owner, name: $repo) {
			issues(first: 10, after: $cursor, states: OPEN) {
				pageInfo {
					endCursor
					hasNextPage
				}
				nodes {
					title
					number
				}
			}
		}
	}`

	vars := map[string]any{
		"owner": "octocat",
		"repo":  "Hello-World",
	}

	result := struct {
		graphql.RateLimit `json:"rateLimit"`

		Repository struct {
			Issues struct {
				graphql.PageInfo `json:"pageInfo"`

				Nodes []struct {
					Title  string `json:"title"`
					Number int    `json:"number"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"repository"`
	}{}

	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)
	g := graphql.New(client)

	if err := g.Query(context.Background(), query, vars, &result); err != nil {
		panic(err)
	}

	fmt.Printf("Result: %+v\n", result)
}
