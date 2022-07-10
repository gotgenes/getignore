package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/google/go-github/v45/github"
	"github.com/gotgenes/getignore/pkg/getignore"
)

// DefaultMaxRequests is the default maximum number of concurrent requests
var DefaultMaxRequests = runtime.NumCPU() - 1

// Getter lists and gets files using the GitHub tree API.
type Getter struct {
	client      *github.Client
	BaseURL     string
	Owner       string
	Repository  string
	Branch      string
	Suffix      string
	MaxRequests int
}

// getterParams holds parameters for instantiating a Getter
type getterParams struct {
	client      *http.Client
	baseURL     string
	owner       string
	repository  string
	branch      string
	suffix      string
	maxRequests int
}

func NewGetter(options ...GetterOption) (Getter, error) {
	params := &getterParams{
		owner:       Owner,
		repository:  Repository,
		branch:      Branch,
		suffix:      Suffix,
		maxRequests: DefaultMaxRequests,
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
	userAgentString := fmt.Sprintf(userAgentTemplate, getignore.Version)
	ghClient.UserAgent = userAgentString
	return Getter{
		client:      ghClient,
		BaseURL:     params.baseURL,
		Owner:       params.owner,
		Repository:  params.repository,
		Branch:      params.branch,
		Suffix:      params.suffix,
		MaxRequests: params.maxRequests,
	}, nil
}

type GetterOption func(*getterParams)

// WithClient sets the HTTP client for the Getter
func WithClient(client *http.Client) GetterOption {
	return func(p *getterParams) {
		p.client = client
	}
}

// WithBaseURL sets the base URL for the Getter
func WithBaseURL(baseURL string) GetterOption {
	return func(p *getterParams) {
		p.baseURL = baseURL
	}
}

// WithOwner sets the owner or organization name for the Getter
func WithOwner(owner string) GetterOption {
	return func(p *getterParams) {
		p.owner = owner
	}
}

// WithRepository sets the repository name for the Getter
func WithRepository(repository string) GetterOption {
	return func(p *getterParams) {
		p.repository = repository
	}
}

// WithBranch sets the branch name for the Getter
func WithBranch(branch string) GetterOption {
	return func(p *getterParams) {
		p.branch = branch
	}
}

// WithSuffix sets the suffix to filter ignore files for
func WithSuffix(suffix string) GetterOption {
	return func(p *getterParams) {
		p.suffix = suffix
	}
}

// WithMaxRequests sets the number of maximum concurrent HTTP requests
func WithMaxRequests(max int) GetterOption {
	return func(p *getterParams) {
		p.maxRequests = max
	}
}

// List returns an array of files filtered by the provided suffix.
func (g Getter) List(ctx context.Context) ([]string, error) {
	tree, err := g.getTree(ctx)
	if err != nil {
		return nil, g.newListError(err)
	}
	entries := g.filterTreeEntries(tree.Entries)
	var files []string
	for _, entry := range entries {
		files = append(files, entry.GetPath())
	}
	return files, nil
}

// Get returns an array of contents of the files downloaded from the given names
func (g Getter) Get(ctx context.Context, names []string) ([]getignore.NamedContents, error) {
	tree, err := g.getTree(ctx)
	if err != nil {
		return nil, g.newGetError(err)
	}
	pathsToSHAs := createPathsToSHAs(tree.Entries)

	names = g.ensureSuffixes(names)
	numNames := len(names)
	namesChan, contentsChan, failedFilesChan := g.startDownloaders(ctx, numNames, pathsToSHAs)

	namesOrdering := createNamesOrdering(names)
	wg, outputChan, errorsChan := startProcessors(namesOrdering, contentsChan, failedFilesChan)

	for _, name := range names {
		namesChan <- name
		wg.Add(1)
	}
	wg.Wait()
	close(namesChan)
	close(contentsChan)
	close(failedFilesChan)

	namedContents := <-outputChan
	failedFiles := <-errorsChan
	if failedFiles != nil {
		err = g.newGetError(failedFiles)
	}
	return namedContents, err
}

func (g Getter) getBlob(
	ctx context.Context,
	pathsToSHAs map[string]string,
	namesChan chan string,
	contentsChan chan getignore.NamedContents,
	failedFilesChan chan getignore.FailedFile,
) {
	for name := range namesChan {
		sha, ok := pathsToSHAs[name]
		if ok {
			blobContents, _, err := g.client.Git.GetBlobRaw(ctx, g.Owner, g.Repository, sha)
			if err != nil {
				failedFile := getignore.FailedFile{
					Name:    name,
					Message: "failed to download",
					Err:     err,
				}
				failedFilesChan <- failedFile
			} else {
				nc := getignore.NamedContents{
					Name:     name,
					Contents: string(blobContents),
				}
				contentsChan <- nc
			}
		} else {
			failedFile := getignore.FailedFile{
				Name:    name,
				Message: "not present in file tree",
			}
			failedFilesChan <- failedFile
		}
	}
}

func (g Getter) newListError(err error) error {
	return fmt.Errorf(
		"error listing contents of %s/%s at %s: %w",
		g.Owner,
		g.Repository,
		g.Branch,
		err,
	)
}

func (g Getter) newGetError(err error) error {
	return fmt.Errorf(
		"error getting files from %s/%s at %s: %w",
		g.Owner,
		g.Repository,
		g.Branch,
		err,
	)
}

func (g Getter) getTree(ctx context.Context) (*github.Tree, error) {
	branch, _, err := g.client.Repositories.GetBranch(ctx, g.Owner, g.Repository, g.Branch, true)
	if err != nil {
		return nil, errors.New("unable to get branch information")
	}
	sha := branch.GetCommit().GetCommit().GetTree().GetSHA()
	if sha == "" {
		return nil, errors.New("no branch information received")
	}
	tree, _, err := g.client.Git.GetTree(ctx, g.Owner, g.Repository, sha, true)
	if err != nil {
		return nil, errors.New("unable to get tree information")
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

func (g Getter) startDownloaders(
	ctx context.Context,
	numFilesToDownload int,
	pathsToSHAs map[string]string,
) (chan string, chan getignore.NamedContents, chan getignore.FailedFile) {
	namesChan := make(chan string, numFilesToDownload)
	maxRequests := min(numFilesToDownload, g.MaxRequests)
	contentsChan := make(chan getignore.NamedContents, numFilesToDownload)
	failedFilesChan := make(chan getignore.FailedFile, numFilesToDownload)
	for i := 0; i < maxRequests; i++ {
		go g.getBlob(ctx, pathsToSHAs, namesChan, contentsChan, failedFilesChan)
	}
	return namesChan, contentsChan, failedFilesChan
}

func (g Getter) ensureSuffixes(names []string) []string {
	paths := make([]string, len(names))
	for i, name := range names {
		path := name
		if filepath.Ext(name) == "" {
			path = name + g.Suffix
		}
		paths[i] = path
	}
	return paths
}

func createPathsToSHAs(entries []*github.TreeEntry) map[string]string {
	pathsToSHAs := make(map[string]string)
	for _, entry := range entries {
		if path := entry.GetPath(); path != "" {
			pathsToSHAs[path] = entry.GetSHA()
		}
	}
	return pathsToSHAs
}

func min(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

func createNamesOrdering(names []string) map[string]int {
	namesOrdering := make(map[string]int)
	for i, name := range names {
		namesOrdering[name] = i
	}
	return namesOrdering
}

func startProcessors(
	namesOrdering map[string]int,
	contentsChan chan getignore.NamedContents,
	failedFilesChan chan getignore.FailedFile,
) (*sync.WaitGroup, chan []getignore.NamedContents, chan getignore.FailedFiles) {
	var wg sync.WaitGroup
	outputChan := make(chan []getignore.NamedContents)
	errorsChan := make(chan getignore.FailedFiles)
	go processContents(contentsChan, namesOrdering, outputChan, &wg)
	go processErrors(failedFilesChan, errorsChan, &wg)
	return &wg, outputChan, errorsChan
}

func processContents(
	contentsChan chan getignore.NamedContents,
	namesOrdering map[string]int,
	outputChannel chan []getignore.NamedContents,
	wg *sync.WaitGroup,
) {
	var allRetrievedContents []getignore.NamedContents
	for contents := range contentsChan {
		allRetrievedContents = append(allRetrievedContents, contents)
		wg.Done()
	}
	sort.Sort(&contentsWithOrdering{contents: allRetrievedContents, ordering: namesOrdering})
	outputChannel <- allRetrievedContents
}

type contentsWithOrdering struct {
	contents []getignore.NamedContents
	ordering map[string]int
}

func (cwo *contentsWithOrdering) Len() int {
	return len(cwo.contents)
}

func (cwo *contentsWithOrdering) Swap(i, j int) {
	cwo.contents[i], cwo.contents[j] = cwo.contents[j], cwo.contents[i]
}

func (cwo *contentsWithOrdering) Less(i, j int) bool {
	return cwo.ordering[cwo.contents[i].Name] < cwo.ordering[cwo.contents[j].Name]
}

func processErrors(
	failedFilesChan chan getignore.FailedFile,
	errorsChan chan getignore.FailedFiles,
	wg *sync.WaitGroup,
) {
	var failedFiles getignore.FailedFiles
	for failedFile := range failedFilesChan {
		failedFiles = append(failedFiles, failedFile)
		wg.Done()
	}
	errorsChan <- failedFiles
}
