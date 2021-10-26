package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/identifiers"
)

// Getter lists ignore files using the GitHub tree API.
type Getter struct {
	client     *github.Client
	BaseURL    string
	Owner      string
	Repository string
	Branch     string
	Suffix     string
}

// gitHubListerParams holds parameters for instantiating a GitHubLister
type gitHubListerParams struct {
	client     *http.Client
	baseURL    string
	owner      string
	repository string
	branch     string
	suffix     string
}

func NewGetter(options ...GitHubListerOption) (Getter, error) {
	params := &gitHubListerParams{
		owner:      Owner,
		repository: Repository,
		branch:     Branch,
		suffix:     Suffix,
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
			return Getter{}, err
		}
	} else {
		ghClient = github.NewClient(params.client)
	}
	userAgentString := fmt.Sprintf(userAgentTemplate, identifiers.Version)
	ghClient.UserAgent = userAgentString
	return Getter{
		client:     ghClient,
		BaseURL:    params.baseURL,
		Owner:      params.owner,
		Repository: params.repository,
		Branch:     params.branch,
		Suffix:     params.suffix,
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

// WithOwner sets the owner or organization name for the GitHubLister
func WithOwner(owner string) GitHubListerOption {
	return func(p *gitHubListerParams) {
		p.owner = owner
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
func (g Getter) List(ctx context.Context) ([]string, error) {
	tree, err := g.getTree(ctx)
	if err != nil {
		return nil, err
	}
	entries := g.filterTreeEntries(tree.Entries)
	var files []string
	for _, entry := range entries {
		files = append(files, entry.GetPath())
	}
	return files, nil
}

func (g Getter) Get(ctx context.Context, names []string) ([]contentstructs.NamedIgnoreContents, error) {
	tree, err := g.getTree(ctx)
	if err != nil {
		return nil, err
	}
	var namedContents []contentstructs.NamedIgnoreContents
	pathsToSHAs := make(map[string]string)
	for _, entry := range tree.Entries {
		if path := entry.GetPath(); path != "" {
			pathsToSHAs[path] = entry.GetSHA()
		}
	}
	for _, name := range names {
		sha, ok := pathsToSHAs[name]
		if ok {
			contents, _, _ := g.client.Git.GetBlobRaw(ctx, g.Owner, g.Repository, sha)
			namedContents = append(namedContents, contentstructs.NamedIgnoreContents{
				Name:     name,
				Contents: string(contents),
			})
		}
	}
	return namedContents, nil
}

func (g Getter) getTree(ctx context.Context) (*github.Tree, error) {
	branch, _, err := g.client.Repositories.GetBranch(ctx, g.Owner, g.Repository, g.Branch, true)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get branch information for %s/%s at %s: %w",
			g.Owner,
			g.Repository,
			g.Branch,
			err,
		)
	}
	sha := branch.GetCommit().GetCommit().GetTree().GetSHA()
	if sha == "" {
		return nil, fmt.Errorf(
			"no branch information received for %s/%s at %s",
			g.Owner,
			g.Repository,
			g.Branch,
		)
	}
	tree, _, err := g.client.Git.GetTree(ctx, g.Owner, g.Repository, sha, true)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to get tree information for %s/%s at %s: %w",
			g.Owner,
			g.Repository,
			g.Branch,
			err,
		)
	}
	return tree, nil
}

func (g Getter) filterTreeEntries(treeEntries []*github.TreeEntry) []*github.TreeEntry {
	var entries []*github.TreeEntry
	for _, entry := range treeEntries {
		if entry.GetType() == "blob" {
			if strings.HasSuffix(entry.GetPath(), g.Suffix) {
				entries = append(entries, entry)
			}
		}
	}
	return entries
}
