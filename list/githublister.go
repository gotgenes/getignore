package list

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v39/github"
	"github.com/gotgenes/getignore/identifiers"
)

// GitHubLister lists ignore files using the GitHub tree API.
type GitHubLister struct {
	client       *github.Client
	BaseURL      string
	Organization string
	Repository   string
	Branch       string
}

// gitHubListerParams holds parameters for instantiating a GitHubLister
type gitHubListerParams struct {
	client       *http.Client
	baseURL      string
	organization string
	repository   string
	branch       string
}

// NewGitHubLister returns a GitHubLister.
func NewGitHubLister(options ...GitHubListerOption) (GitHubLister, error) {
	params := &gitHubListerParams{
		baseURL:      "https://api.github.com",
		organization: "github",
		repository:   "gitignore",
		branch:       "master",
	}
	for _, option := range options {
		option(params)
	}
	var (
		ghClient *github.Client
		err      error
	)
	if params.baseURL != "" {
		ghClient, err = github.NewEnterpriseClient(params.baseURL, params.baseURL, params.client)
		if err != nil {
			return GitHubLister{}, err
		}
	} else {
		ghClient = github.NewClient(params.client)
	}
	userAgentString := fmt.Sprintf(userAgentTemplate, identifiers.Version)
	ghClient.UserAgent = userAgentString
	return GitHubLister{
		client:       ghClient,
		BaseURL:      params.baseURL,
		Organization: params.organization,
		Repository:   params.repository,
		Branch:       params.branch,
	}, nil
}

type GitHubListerOption func(*gitHubListerParams)

// WithClient sets the HTTP client for the GitHubLister
func WithClient(client *http.Client) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.client = client
	}
}

// WithBaseURL sets the base URL for the GitHubLister
func WithBaseURL(baseURL string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.baseURL = baseURL
	}
}

// WithOrganization sets the organization or user name for the GitHubLister
func WithOrganization(organization string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.organization = organization
	}
}

// WithRepository sets the repository name for the GitHubLister
func WithRepository(repository string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.repository = repository
	}
}

// WithBranch sets the branch name for the GitHubLister
func WithBranch(branch string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.branch = branch
	}
}

// List returns an array of ignore files filtered by the provided suffix.
// Passing an empty string for suffix will return all files, with no filtering.
func (l GitHubLister) List(ctx context.Context, suffix string) ([]string, error) {
	var files []string
	branch, _, _ := l.client.Repositories.GetBranch(ctx, l.Organization, l.Repository, l.Branch, true)
	sha := branch.GetCommit().GetCommit().GetTree().GetSHA()
	l.client.Git.GetTree(ctx, l.Organization, l.Repository, sha, true)
	return files, nil
}
