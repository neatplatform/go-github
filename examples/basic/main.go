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
