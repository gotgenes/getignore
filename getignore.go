package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

func getNamesFromArguments(context *cli.Context) []string {
	names := context.Args()

	if context.String("names-file") != "" {
		namesFile, _ := os.Open(context.String("names-file"))
		names = append(names, ParseNamesFile(namesFile)...)
	}
	return names
}

// ParseNamesFile reads a file containing one name of a gitignore patterns file per line
func ParseNamesFile(namesFile io.Reader) []string {
	var a []string
	scanner := bufio.NewScanner(namesFile)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if len(name) > 0 {
			a = append(a, name)
		}
	}
	return a
}

// CreateNamesOrdering creates a mapping of each name to its respective input position
func CreateNamesOrdering(names []string) map[string]int {
	namesOrdering := make(map[string]int)
	for i, name := range names {
		namesOrdering[name] = i
	}
	return namesOrdering
}

// HTTPIgnoreGetter provides an implementation to retrieve gitignore patterns from
// files available over HTTP
type HTTPIgnoreGetter struct {
	baseURL          string
	defaultExtension string
	maxConnections   int
}

// RetrievedContents represents the result of retrieving contents of a gitignore patterns
// file
type RetrievedContents struct {
	name     string
	source   string
	contents string
	err      error
}

// GetIgnoreFiles retrieves gitignore patterns files via HTTP and sends their contents
// over a channel. It registers each request made with a WaitGroup instance, so the
// responses can be awaited.
func (getter *HTTPIgnoreGetter) GetIgnoreFiles(names []string, contentsChannel chan RetrievedContents, requestsPending *sync.WaitGroup) {
	namesChannel := make(chan string)
	for i := 0; i < getter.maxConnections; i++ {
		go getter.downloadIgnoreFile(namesChannel, contentsChannel, requestsPending)
	}
	for _, name := range names {
		requestsPending.Add(1)
		namesChannel <- name
	}
	close(namesChannel)
}

// FailedSource represents a source unable to be retrieved or processed
type FailedSource struct {
	source string
	err    error
}

func (fs *FailedSource) Error() string {
	return fmt.Sprintf("%s %s", fs.source, fs.err.Error())
}

func (getter *HTTPIgnoreGetter) downloadIgnoreFile(namesChannel chan string, contentsChannel chan RetrievedContents, requestsPending *sync.WaitGroup) {
	for name := range namesChannel {
		url := getter.nameToURL(name)
		log.Println("Retrieving", url)
		response, err := http.Get(url)
		contents, err := getter.processResponse(response, err)
		contentsChannel <- RetrievedContents{name, url, contents, err}
		requestsPending.Done()
	}
}

func (getter *HTTPIgnoreGetter) nameToURL(name string) string {
	nameWithExtension := getter.getNameWithExtension(name)
	url := getter.baseURL + "/" + nameWithExtension
	return url
}

func (getter *HTTPIgnoreGetter) getNameWithExtension(name string) string {
	if filepath.Ext(name) == "" {
		name = name + getter.defaultExtension
	}
	return name
}

func getContent(body io.ReadCloser) (content string, err error) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		content = content + fmt.Sprintln(scanner.Text())
	}
	err = scanner.Err()
	return content, err
}

func (getter *HTTPIgnoreGetter) processResponse(response *http.Response, err error) (contents string, processedErr error) {
	if err != nil {
		processedErr = err
	} else if response.StatusCode != 200 {
		processedErr = fmt.Errorf("Got status code %d", response.StatusCode)
	} else {
		defer response.Body.Close()
		var contentErr error
		contents, contentErr = getContent(response.Body)
		if contentErr != nil {
			processedErr = fmt.Errorf("Error reading response body: %s", contentErr.Error())
		}
	}
	return
}

// FailedSources represents a collection of FailedSource instances
type FailedSources struct {
	sources []*FailedSource
}

// Add adds a FailedSource instance to the FailedSources collection
func (failedSources *FailedSources) Add(failedSource *FailedSource) {
	failedSources.sources = append(failedSources.sources, failedSource)
}

func (failedSources *FailedSources) Error() string {
	sourceErrors := make([]string, len(failedSources.sources))
	for i, failedSource := range failedSources.sources {
		sourceErrors[i] = failedSource.Error()
	}
	stringOfErrors := strings.Join(sourceErrors, "\n")
	return "Errors for the following URLs:\n" + stringOfErrors
}

// NamedIgnoreContents represents the contents (patterns and comments) of a
// gitignore file
type NamedIgnoreContents struct {
	name     string
	contents string
}

// DisplayName returns the decorated name, suitable for a section header in a
// gitignore file
func (nic *NamedIgnoreContents) DisplayName() string {
	baseName := filepath.Base(nic.name)
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}

func processContents(contentsChannel chan RetrievedContents, namesOrdering map[string]int) ([]NamedIgnoreContents, error) {
	allRetrievedContents := make([]NamedIgnoreContents, len(namesOrdering))
	var err error
	failedSources := new(FailedSources)
	for retrievedContents := range contentsChannel {
		if retrievedContents.err != nil {
			failedSource := &FailedSource{retrievedContents.source, retrievedContents.err}
			failedSources.Add(failedSource)
		} else {
			name := retrievedContents.name
			position, present := namesOrdering[name]
			if !present {
				return allRetrievedContents, fmt.Errorf("Could not find name %s in ordering", name)
			}
			allRetrievedContents[position] = NamedIgnoreContents{name, retrievedContents.contents}
		}
	}
	if len(failedSources.sources) > 0 {
		err = failedSources
	}
	return allRetrievedContents, err
}

func getOutputFile(context *cli.Context) (outputFilePath string, outputFile io.Writer, err error) {
	outputFilePath = context.String("o")
	if outputFilePath == "" {
		outputFilePath = "STDOUT"
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(outputFilePath)
	}
	return
}

func writeIgnoreFile(ignoreFile io.Writer, contents []NamedIgnoreContents) (err error) {
	writer := bufio.NewWriter(ignoreFile)
	for i, nc := range contents {
		if i > 0 {
			writer.WriteString("\n\n")
		}
		writer.WriteString(decorateName(nc.DisplayName()))
		writer.WriteString(nc.contents)
	}
	if writer.Flush() != nil {
		err = writer.Flush()
	}
	return
}

func decorateName(name string) string {
	nameLength := len(name)
	fullHashLine := strings.Repeat("#", nameLength+4)
	nameLine := fmt.Sprintf("# %s #", name)
	decoratedName := strings.Join([]string{fullHashLine, nameLine, fullHashLine, ""}, "\n")
	return decoratedName
}

func creatCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "getignore"
	app.Version = "0.2.0.dev"
	app.Usage = "Bootstraps gitignore files from central sources"

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "get",
			Usage: "Retrieves gitignore patterns files from a central source and concatenates them",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "base-url, u",
					Usage: "The URL under which gitignore files can be found",
					Value: "https://raw.githubusercontent.com/github/gitignore/master",
				},
				cli.StringFlag{
					Name:  "default-extension, e",
					Usage: "The default file extension appended to names when retrieving them",
					Value: ".gitignore",
				},
				cli.IntFlag{
					Name:  "max-connections",
					Usage: "The number of maximum connections to open for HTTP requests",
					Value: 8,
				},
				cli.StringFlag{
					Name:  "names-file, n",
					Usage: "Path to file containing names of gitignore patterns files",
				},
				cli.StringFlag{
					Name:  "o",
					Usage: "Path to output file (default: STDOUT)",
				},
			},
			ArgsUsage: "[gitignore_name] [gitignore_name …]",
			Action:    downloadAllIgnoreFiles,
		},
	}

	return app
}

func downloadAllIgnoreFiles(context *cli.Context) error {
	names := getNamesFromArguments(context)
	namesOrdering := CreateNamesOrdering(names)
	getter := HTTPIgnoreGetter{context.String("base-url"), context.String("default-extension"), context.Int("max-connections")}
	contentsChannel := make(chan RetrievedContents, context.Int("max-connections"))
	var requestsPending sync.WaitGroup
	getter.GetIgnoreFiles(names, contentsChannel, &requestsPending)
	requestsPending.Wait()
	close(contentsChannel)
	contents, err := processContents(contentsChannel, namesOrdering)
	if err != nil {
		return err
	}
	outputFilePath, outputFile, err := getOutputFile(context)
	if err != nil {
		return err
	}
	log.Println("Writing contents to", outputFilePath)
	err = writeIgnoreFile(outputFile, contents)
	if err != nil {
		return err
	}
	return err
}

func main() {
	log.SetFlags(0)
	app := creatCLI()
	app.RunAndExitOnError()
}
