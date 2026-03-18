package rest_test

import (
	"context"
	"fmt"

	"github.com/neatplatform/go-github/rest"
)

func ExampleClient_EnsureScopes() {
	client := rest.NewClient("")
	if err := client.EnsureScopes(context.Background(), rest.ScopeRepo); err != nil {
		panic(err)
	}
}

func ExampleUserService_Get() {
	client := rest.NewClient("")
	user, resp, err := client.Users.Get(context.Background(), "octocat")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rate: %+v\n\n", resp.Rate)
	fmt.Printf("Name: %s\n", user.Name)
}

func ExampleRepoService_Commits() {
	client := rest.NewClient("")
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
	client := rest.NewClient("")
	issues, resp, err := client.Repo("octocat", "Hello-World").Issues.List(context.Background(), 50, 1, rest.IssuesFilter{})
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
	client := rest.NewClient("")
	pulls, resp, err := client.Repo("octocat", "Hello-World").Pulls.List(context.Background(), 50, 1, rest.PullsFilter{})
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
	client := rest.NewClient("")
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
	client := rest.NewClient("")

	query := rest.SearchQuery{}
	query.IncludeKeywords("Fix")
	query.ExcludeKeywords("WIP")
	query.IncludeQualifiers(
		rest.QualifierTypePR,
		rest.QualifierInTitle,
		rest.QualifierLabel("bug"),
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
