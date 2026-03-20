package github

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
)

// ReleaseService provides GitHub REST APIs for releases in a repository.
//
// See https://docs.github.com/en/rest/releases/releases
type ReleaseService struct {
	client      *Client
	owner, repo string
}

// List retrieves all releases.
// If the user has push access, draft releases will also be returned.
//
// See https://docs.github.com/en/rest/releases/releases#list-releases
func (s *ReleaseService) List(ctx context.Context, pageSize, pageNo int) ([]Release, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases", s.owner, s.repo)

	req, err := s.client.NewPageRequest(ctx, "GET", url, pageSize, pageNo, nil)
	if err != nil {
		return nil, nil, err
	}

	releases := []Release{}
	resp, err := s.client.Do(req, &releases)
	if err != nil {
		return nil, nil, err
	}

	return releases, resp, nil
}

// Latest returns the latest release.
// The latest release is the most recent non-prerelease and non-draft release.
//
// See https://docs.github.com/en/rest/releases/releases#get-the-latest-release
func (s *ReleaseService) Latest(ctx context.Context) (*Release, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases/latest", s.owner, s.repo)

	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	release := new(Release)

	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, nil, err
	}

	return release, resp, nil
}

// Get retrieves a release by release id.
//
// See https://docs.github.com/en/rest/releases/releases#get-a-release
func (s *ReleaseService) Get(ctx context.Context, id int) (*Release, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases/%d", s.owner, s.repo, id)

	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	release := new(Release)

	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, nil, err
	}

	return release, resp, nil
}

// GetByTag retrieves a release by tag name.
//
// See https://docs.github.com/en/rest/releases/releases#get-a-release-by-tag-name
func (s *ReleaseService) GetByTag(ctx context.Context, tag string) (*Release, *Response, error) {
	escapedTag := url.PathEscape(tag)
	url := fmt.Sprintf("/repos/%s/%s/releases/tags/%s", s.owner, s.repo, escapedTag)

	req, err := s.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	release := new(Release)

	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, nil, err
	}

	return release, resp, nil
}

// Create creates a new GitHub release.
//
// See https://docs.github.com/en/rest/releases/releases#create-a-release
func (s *ReleaseService) Create(ctx context.Context, params ReleaseParams) (*Release, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases", s.owner, s.repo)

	req, err := s.client.NewRequest(ctx, "POST", url, params)
	if err != nil {
		return nil, nil, err
	}

	release := new(Release)

	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, nil, err
	}

	return release, resp, nil
}

// Update updates an existing GitHub release.
//
// See https://docs.github.com/en/rest/releases/releases#update-a-release
func (s *ReleaseService) Update(ctx context.Context, id int, params ReleaseParams) (*Release, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases/%d", s.owner, s.repo, id)

	req, err := s.client.NewRequest(ctx, "PATCH", url, params)
	if err != nil {
		return nil, nil, err
	}

	release := new(Release)

	resp, err := s.client.Do(req, release)
	if err != nil {
		return nil, nil, err
	}

	return release, resp, nil
}

// Delete deletes a release by release id.
//
// See https://docs.github.com/en/rest/releases/releases#delete-a-release
func (s *ReleaseService) Delete(ctx context.Context, id int) (*Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases/%d", s.owner, s.repo, id)

	req, err := s.client.NewRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UploadAsset uploads a file to a GitHub release.
//
// See https://docs.github.com/en/rest/releases/assets#upload-a-release-asset
func (s *ReleaseService) UploadAsset(ctx context.Context, id int, assetFile, assetLabel string) (*ReleaseAsset, *Response, error) {
	url := fmt.Sprintf("/repos/%s/%s/releases/%d/assets", s.owner, s.repo, id)

	req, closer, err := s.client.NewUploadRequest(ctx, url, assetFile)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = closer.Close()
	}()

	q := req.URL.Query()
	if assetName := filepath.Base(assetFile); assetName != "" {
		q.Add("name", assetName)
	}
	if assetLabel != "" {
		q.Add("label", assetLabel)
	}
	req.URL.RawQuery = q.Encode()

	asset := new(ReleaseAsset)
	resp, err := s.client.Do(req, asset)
	if err != nil {
		return nil, nil, err
	}

	return asset, resp, nil
}

// DownloadAsset downloads an asset from a GitHub release.
//
// See https://docs.github.com/en/rest/releases/assets#get-a-release-asset
func (s *ReleaseService) DownloadAsset(ctx context.Context, tag, assetName string, w io.Writer) (*Response, error) {
	escapedTag := url.PathEscape(tag)
	escapedAssetName := url.PathEscape(assetName)
	url := fmt.Sprintf("/%s/%s/releases/download/%s/%s", s.owner, s.repo, escapedTag, escapedAssetName)

	req, err := s.client.NewDownloadRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, w)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
