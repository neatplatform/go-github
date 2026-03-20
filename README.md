[![Go Doc][godoc-image]][godoc-url]
[![Build Status][workflow-image]][workflow-url]
[![Test Coverage][codecov-image]][codecov-url]

# go-github

A simple Go client for [GitHub REST API](https://docs.github.com/en/rest) and [GitHub GraphQL API](https://docs.github.com/en/graphql).

## Quick Start

### REST API

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neatplatform/go-github"
)

func main() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	commits, resp, err := client.Repo("octocat", "Hello-World").Commits(context.Background(), 50, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, commit := range commits {
		fmt.Printf("%s\n", commit.SHA)
	}
}
```

### GraphQL API

```go
package main

import (
  "context"
  "fmt"
  "os"

  "github.com/neatplatform/go-github"
  "github.com/neatplatform/go-github/graphql"
)

func main() {
  const query = `query {
    viewer {
      login
    }
  }`

  result := struct {
    Viewer struct {
      Login string `json:"login"`
    } `json:"viewer"`
  }{}

	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)
  g := graphql.New(client)

  if err := g.Query(context.Background(), query, nil, &result); err != nil {
    panic(err)
  }

  fmt.Printf("Result: %+v\n", result)
}
```

## Resources

  - [Best practices for using the REST API](https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api)


[godoc-url]: https://pkg.go.dev/github.com/neatplatform/go-github
[godoc-image]: https://pkg.go.dev/badge/github.com/neatplatform/go-github
[workflow-url]: https://github.com/neatplatform/go-github/actions/workflows/go.yml
[workflow-image]: https://github.com/neatplatform/go-github/actions/workflows/go.yml/badge.svg
[codecov-url]: https://codecov.io/gh/neatplatform/go-github
[codecov-image]: https://codecov.io/gh/neatplatform/go-github/graph/badge.svg
