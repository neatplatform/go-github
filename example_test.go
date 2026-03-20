package github_test

import (
	"context"
	"fmt"
	"os"

	"github.com/neatplatform/go-github"
)

func ExampleClient_EnsureScopes() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	if err := client.EnsureScopes(context.Background(), github.ScopeRepo); err != nil {
		panic(err)
	}
}

func ExampleUserService_Get() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	user, resp, err := client.Users.Get(context.Background(), "octocat")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	fmt.Printf("Name: %s\n", user.Name)
}

func ExampleRepoService_Commits() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	commits, resp, err := client.Repo("octocat", "Hello-World").Commits(context.Background(), 50, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, c := range commits {
		fmt.Printf("SHA: %s\n", c.SHA)
	}
}

func ExampleIssueService_List() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	issues, resp, err := client.Repo("octocat", "Hello-World").Issues.List(context.Background(), 50, 1, github.IssuesFilter{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, i := range issues {
		fmt.Printf("Title: %s\n", i.Title)
	}
}

func ExamplePullService_List() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	pulls, resp, err := client.Repo("octocat", "Hello-World").Pulls.List(context.Background(), 50, 1, github.PullsFilter{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, p := range pulls {
		fmt.Printf("Title: %s\n", p.Title)
	}
}

func ExampleReleaseService_List() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	releases, resp, err := client.Repo("octocat", "Hello-World").Releases.List(context.Background(), 20, 1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, r := range releases {
		fmt.Printf("Name: %s\n", r.Name)
	}
}

func ExampleSearchService_SearchIssues() {
	authToken := os.Getenv("GITHUB_TOKEN")
	client := github.NewClient(authToken)

	query := github.SearchQuery{}
	query.IncludeKeywords("Fix")
	query.ExcludeKeywords("WIP")
	query.IncludeQualifiers(
		github.QualifierTypePR,
		github.QualifierInTitle,
		github.QualifierLabel("bug"),
	)

	result, resp, err := client.Search.SearchIssues(context.Background(), 20, 1, "", "", query)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pages: %+v\n", resp.Pages)
	fmt.Printf("Rate: %+v\n\n", resp.Rate)

	for _, issue := range result.Items {
		fmt.Printf("%s\n", issue.HTMLURL)
	}
}
