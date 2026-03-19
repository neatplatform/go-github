package github

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const dateYYYYMMDD = "2006-01-02"

// SearchResultSort determines how to sort search results.
type SearchResultSort string

const (
	SortByDefault      SearchResultSort = ""
	SortByFollowers    SearchResultSort = "followers"
	SortByRepositories SearchResultSort = "repositories"
	SortByJoined       SearchResultSort = "joined"
	SortByStars        SearchResultSort = "stars"
	SortByForks        SearchResultSort = "forks"
	SortByUpdated      SearchResultSort = "updated"
	SortByCreated      SearchResultSort = "created"
	SortByComments     SearchResultSort = "comments"
	SortByReactions    SearchResultSort = "reactions"
	SortByInteractions SearchResultSort = "interactions"
)

// SearchResultOrder determines the order of search results.
type SearchResultOrder string

const (
	DefaultOrder SearchResultOrder = ""
	AscOrder     SearchResultOrder = "asc"
	DescOrder    SearchResultOrder = "desc"
)

// Qualifier is a GitHub search qualifier.
//
// See https://docs.github.com/en/search-github/searching-on-github
type Qualifier string

const (
	QualifierTypeUser      Qualifier = "type:user"
	QualifierTypeOrg       Qualifier = "type:org"
	QualifierTypePR        Qualifier = "type:pr"
	QualifierTypeIssue     Qualifier = "type:issue"
	QualifierIsPR          Qualifier = "is:pr"
	QualifierIsIssue       Qualifier = "is:issue"
	QualifierInLogin       Qualifier = "in:login"
	QualifierInEmail       Qualifier = "in:email"
	QualifierInName        Qualifier = "in:name"
	QualifierInDescription Qualifier = "in:description"
	QualifierInREADME      Qualifier = "in:readme"
	QualifierInTitle       Qualifier = "in:title"
	QualifierInBody        Qualifier = "in:body"
	QualifierInComments    Qualifier = "in:comments"
	QualifierStateOpen     Qualifier = "state:open"
	QualifierStateClosed   Qualifier = "state:closed"
	QualifierIsOpen        Qualifier = "is:open"
	QualifierIsClosed      Qualifier = "is:closed"
	QualifierIsPublic      Qualifier = "is:public"
	QualifierIsInternal    Qualifier = "is:internal"
	QualifierIsPrivate     Qualifier = "is:private"
	QualifierStatusPending Qualifier = "status:pending"
	QualifierStatusSuccess Qualifier = "status:success"
	QualifierStatusFailure Qualifier = "status:failure"
	QualifierArchivedTrue  Qualifier = "archived:true"
	QualifierArchivedFalse Qualifier = "archived:false"
	QualifierDraftTrue     Qualifier = "draft:true"
	QualifierDraftFalse    Qualifier = "draft:false"
	QualifierIsMerged      Qualifier = "is:merged"
	QualifierIsUnmerged    Qualifier = "is:unmerged"
	QualifierIsLocked      Qualifier = "is:locked"
	QualifierIsUnlocked    Qualifier = "is:unlocked"
	QualifierNoAssignee    Qualifier = "no:assignee"
	QualifierNoLabel       Qualifier = "no:label"
	QualifierNoProject     Qualifier = "no:project"
	QualifierNoMilestone   Qualifier = "no:milestone"
)

// QualifierUser creates a qualifier for matching a user.
func QualifierUser(username string) Qualifier {
	return Qualifier(fmt.Sprintf("user:%s", username))
}

// QualifierOrg creates a qualifier for matching an organization.
func QualifierOrg(orgname string) Qualifier {
	return Qualifier(fmt.Sprintf("org:%s", orgname))
}

// QualifierRepo creates a qualifier for matching a repository.
func QualifierRepo(repoOwner, repoName string) Qualifier {
	return Qualifier(fmt.Sprintf("repo:%s/%s", repoOwner, repoName))
}

// QualifierAuthor creates a qualifier for matching an author.
func QualifierAuthor(username string) Qualifier {
	return Qualifier(fmt.Sprintf("author:%s", username))
}

// QualifierAuthorApp creates a qualifier for matching an author app.
func QualifierAuthorApp(username string) Qualifier {
	return Qualifier(fmt.Sprintf("author:app/%s", username))
}

// QualifierAssignee creates a qualifier for matching an assignee.
func QualifierAssignee(username string) Qualifier {
	return Qualifier(fmt.Sprintf("assignee:%s", username))
}

// QualifierLabel creates a qualifier for matching a label.
func QualifierLabel(label string) Qualifier {
	return Qualifier(fmt.Sprintf("label:%q", label))
}

// QualifierMilestone creates a qualifier for matching a milestone.
func QualifierMilestone(milestone string) Qualifier {
	return Qualifier(fmt.Sprintf("milestone:%q", milestone))
}

// QualifierProject creates a qualifier for matching a project board.
func QualifierProject(projectBoard string) Qualifier {
	return Qualifier(fmt.Sprintf("project:%q", projectBoard))
}

// QualifierRepoProject creates a qualifier for matching a repository project board.
func QualifierRepoProject(repoOwner, repoName, projectBoard string) Qualifier {
	return Qualifier(fmt.Sprintf("project:%s/%s/%s", repoOwner, repoName, projectBoard))
}

// QualifierHead creates a qualifier for matching a head branch.
func QualifierHead(branch string) Qualifier {
	return Qualifier(fmt.Sprintf("head:%s", branch))
}

// QualifierBase creates a qualifier for matching a base branch.
func QualifierBase(branch string) Qualifier {
	return Qualifier(fmt.Sprintf("base:%s", branch))
}

// QualifierLanguage creates a qualifier for matching a language.
func QualifierLanguage(language string) Qualifier {
	return Qualifier(fmt.Sprintf("language:%s", language))
}

// QualifierTopic creates a qualifier for matching a topic.
func QualifierTopic(topic string) Qualifier {
	return Qualifier(fmt.Sprintf("topic:%s", topic))
}

// QualifierCreatedOn creates a qualifier for matching issues and pull requests created on a date.
func QualifierCreatedOn(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:%s", t.Format(dateYYYYMMDD)))
}

// QualifierCreatedAfter creates a qualifier for matching issues and pull requests created after a date.
func QualifierCreatedAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:>%s", t.Format(dateYYYYMMDD)))
}

// QualifierCreatedOnOrAfter creates a qualifier for matching issues and pull requests created on or after a date.
func QualifierCreatedOnOrAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:>=%s", t.Format(dateYYYYMMDD)))
}

// QualifierCreatedBefore creates a qualifier for matching issues and pull requests created before a date.
func QualifierCreatedBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:<%s", t.Format(dateYYYYMMDD)))
}

// QualifierCreatedOnOrBefore creates a qualifier for matching issues and pull requests created on or before a date.
func QualifierCreatedOnOrBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:<=%s", t.Format(dateYYYYMMDD)))
}

// QualifierCreatedBetween creates a qualifier for matching issues and pull requests created between two dates.
func QualifierCreatedBetween(from, to time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("created:%s..%s", from.Format(dateYYYYMMDD), to.Format(dateYYYYMMDD)))
}

// QualifierUpdatedOn creates a qualifier for matching issues and pull requests updated on a date.
func QualifierUpdatedOn(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:%s", t.Format(dateYYYYMMDD)))
}

// QualifierUpdatedAfter creates a qualifier for matching issues and pull requests updated after a date.
func QualifierUpdatedAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:>%s", t.Format(dateYYYYMMDD)))
}

// QualifierUpdatedOnOrAfter creates a qualifier for matching issues and pull requests updated on or after a date.
func QualifierUpdatedOnOrAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:>=%s", t.Format(dateYYYYMMDD)))
}

// QualifierUpdatedBefore creates a qualifier for matching issues and pull requests updated before a date.
func QualifierUpdatedBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:<%s", t.Format(dateYYYYMMDD)))
}

// QualifierUpdatedOnOrBefore creates a qualifier for matching issues and pull requests updated on or before a date.
func QualifierUpdatedOnOrBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:<=%s", t.Format(dateYYYYMMDD)))
}

// QualifierUpdatedBetween creates a qualifier for matching issues and pull requests updated between two dates.
func QualifierUpdatedBetween(from, to time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("updated:%s..%s", from.Format(dateYYYYMMDD), to.Format(dateYYYYMMDD)))
}

// QualifierClosedOn creates a qualifier for matching issues and pull requests closed on a date.
func QualifierClosedOn(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:%s", t.Format(dateYYYYMMDD)))
}

// QualifierClosedAfter creates a qualifier for matching issues and pull requests closed after a date.
func QualifierClosedAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:>%s", t.Format(dateYYYYMMDD)))
}

// QualifierClosedOnOrAfter creates a qualifier for matching issues and pull requests closed on or after a date.
func QualifierClosedOnOrAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:>=%s", t.Format(dateYYYYMMDD)))
}

// QualifierClosedBefore creates a qualifier for matching issues and pull requests closed before a date.
func QualifierClosedBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:<%s", t.Format(dateYYYYMMDD)))
}

// QualifierClosedOnOrBefore creates a qualifier for matching issues and pull requests closed on or before a date.
func QualifierClosedOnOrBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:<=%s", t.Format(dateYYYYMMDD)))
}

// QualifierClosedBetween creates a qualifier for matching issues and pull requests closed between two dates.
func QualifierClosedBetween(from, to time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("closed:%s..%s", from.Format(dateYYYYMMDD), to.Format(dateYYYYMMDD)))
}

// QualifierMergedOn creates a qualifier for matching pull requests merged on a date.
func QualifierMergedOn(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:%s", t.Format(dateYYYYMMDD)))
}

// QualifierMergedAfter creates a qualifier for matching pull requests merged after a date.
func QualifierMergedAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:>%s", t.Format(dateYYYYMMDD)))
}

// QualifierMergedOnOrAfter creates a qualifier for matching pull requests merged on or after a date.
func QualifierMergedOnOrAfter(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:>=%s", t.Format(dateYYYYMMDD)))
}

// QualifierMergedBefore creates a qualifier for matching pull requests merged before a date.
func QualifierMergedBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:<%s", t.Format(dateYYYYMMDD)))
}

// QualifierMergedOnOrBefore creates a qualifier for matching pull requests merged on or before a date.
func QualifierMergedOnOrBefore(t time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:<=%s", t.Format(dateYYYYMMDD)))
}

// QualifierMergedBetween creates a qualifier for matching pull requests merged between two dates.
func QualifierMergedBetween(from, to time.Time) Qualifier {
	return Qualifier(fmt.Sprintf("merged:%s..%s", from.Format(dateYYYYMMDD), to.Format(dateYYYYMMDD)))
}

// SearchQuery is used for searching GitHub.
//
// See https://docs.github.com/en/rest/search/search#constructing-a-search-query
// See https://docs.github.com/en/search-github/getting-started-with-searching-on-github/understanding-the-search-syntax
type SearchQuery struct {
	keywords   []string
	qualifiers []Qualifier
}

func (q *SearchQuery) IncludeKeywords(keywords ...string) {
	for _, keyword := range keywords {
		q.keywords = append(q.keywords, fmt.Sprintf("%q", keyword))
	}
}

func (q *SearchQuery) ExcludeKeywords(keywords ...string) {
	for _, keyword := range keywords {
		q.keywords = append(q.keywords, fmt.Sprintf("NOT %q", keyword))
	}
}

func (q *SearchQuery) IncludeQualifiers(qualifiers ...Qualifier) {
	q.qualifiers = append(q.qualifiers, qualifiers...)
}

func (q *SearchQuery) ExcludeQualifiers(qualifiers ...Qualifier) {
	for _, qualifier := range qualifiers {
		q.qualifiers = append(q.qualifiers, Qualifier("-"+qualifier))
	}
}

// String returns the GitHub search query parameter.
func (q *SearchQuery) String() string {
	terms := append([]string{}, q.keywords...)
	for _, qualifier := range q.qualifiers {
		terms = append(terms, string(qualifier))
	}

	return strings.Join(terms, " ")
}

// SearchService provides GitHub REST APIs for searching users, repositories, issues, pull requests, etc.
//
// See https://docs.github.com/en/rest/search/search
type SearchService struct {
	client *Client
}

// SearchUsersResult is the result of searching users.
type SearchUsersResult struct {
	TotalCount        int    `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []User `json:"items"`
}

// SearchUsers searches users.
//
// See https://docs.github.com/en/rest/search/search#search-users
// See https://docs.github.com/en/search-github/searching-on-github/searching-users
func (s *SearchService) SearchUsers(ctx context.Context, pageSize, pageNo int, sort SearchResultSort, order SearchResultOrder, query SearchQuery) (*SearchUsersResult, *Response, error) {
	req, err := s.client.NewPageRequest(ctx, "GET", "/search/users", pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	if sort != "" {
		q.Add("sort", string(sort))
	}
	if order != "" {
		q.Add("order", string(order))
	}
	q.Add("q", query.String())
	req.URL.RawQuery = q.Encode()

	searchResult := new(SearchUsersResult)
	resp, err := s.client.Do(req, searchResult)
	if err != nil {
		return nil, nil, err
	}

	return searchResult, resp, nil
}

// SearchReposResult is the result of searching repositories.
type SearchReposResult struct {
	TotalCount        int          `json:"total_count"`
	IncompleteResults bool         `json:"incomplete_results"`
	Items             []Repository `json:"items"`
}

// SearchRepos searches repositories.
//
// See https://docs.github.com/en/rest/search/search#search-repositories
// See https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories
func (s *SearchService) SearchRepos(ctx context.Context, pageSize, pageNo int, sort SearchResultSort, order SearchResultOrder, query SearchQuery) (*SearchReposResult, *Response, error) {
	req, err := s.client.NewPageRequest(ctx, "GET", "/search/repositories", pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	if sort != "" {
		q.Add("sort", string(sort))
	}
	if order != "" {
		q.Add("order", string(order))
	}
	q.Add("q", query.String())
	req.URL.RawQuery = q.Encode()

	searchResult := new(SearchReposResult)
	resp, err := s.client.Do(req, searchResult)
	if err != nil {
		return nil, nil, err
	}

	return searchResult, resp, nil
}

// SearchIssuesResult is the result of searching issues and pull requests.
type SearchIssuesResult struct {
	TotalCount        int     `json:"total_count"`
	IncompleteResults bool    `json:"incomplete_results"`
	Items             []Issue `json:"items"`
}

// SearchIssues searches issues and pull requests.
//
// See https://docs.github.com/en/rest/search/search#search-issues-and-pull-requests
// See https://docs.github.com/en/search-github/searching-on-github/searching-issues-and-pull-requests
func (s *SearchService) SearchIssues(ctx context.Context, pageSize, pageNo int, sort SearchResultSort, order SearchResultOrder, query SearchQuery) (*SearchIssuesResult, *Response, error) {
	req, err := s.client.NewPageRequest(ctx, "GET", "/search/issues", pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	if sort != "" {
		q.Add("sort", string(sort))
	}
	if order != "" {
		q.Add("order", string(order))
	}
	q.Add("q", query.String())
	req.URL.RawQuery = q.Encode()

	searchResult := new(SearchIssuesResult)
	resp, err := s.client.Do(req, searchResult)
	if err != nil {
		return nil, nil, err
	}

	return searchResult, resp, nil
}
