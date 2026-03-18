package main

import (
	"context"
	"fmt"

	"github.com/neatplatform/go-github"
)

func main() {
	client := github.NewClient("")

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
