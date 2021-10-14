package list

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	Suffix       string
}

// gitHubListerParams holds parameters for instantiating a GitHubLister
type gitHubListerParams struct {
	client       *http.Client
	baseURL      string
	organization string
	repository   string
	branch       string
	suffix       string
}

// NewGitHubLister returns a GitHubLister.
func NewGitHubLister(options ...GitHubListerOption) (GitHubLister, error) {
	params := &gitHubListerParams{
		organization: "github",
		repository:   "gitignore",
		branch:       "master",
		suffix:       ".gitignore",
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
		Suffix:       params.suffix,
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

// WithSuffix sets the suffix to filter ignore files for
func WithSuffix(suffix string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.suffix = suffix
	}
}

// List returns an array of files filtered by the provided suffix.
func (l GitHubLister) List(ctx context.Context) ([]string, error) {
	var files []string
	branch, _, err := l.client.Repositories.GetBranch(ctx, l.Organization, l.Repository, l.Branch, true)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get branch information for %s/%s at %s: %w",
			l.Organization,
			l.Repository,
			l.Branch,
			err,
		)
	}
	sha := branch.GetCommit().GetCommit().GetTree().GetSHA()
	tree, _, _ := l.client.Git.GetTree(ctx, l.Organization, l.Repository, sha, true)
	for _, entry := range tree.Entries {
		if entry.GetType() == "blob" {
			path := entry.GetPath()
			if strings.HasSuffix(path, l.Suffix) {
				files = append(files, path)
			}
		}
	}
	return files, nil
}
