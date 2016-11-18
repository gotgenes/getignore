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
		names = append(names, parseNamesFile(namesFile)...)
	}
	return names
}

func parseNamesFile(namesFile io.Reader) []string {
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

func createNamesOrdering(names []string) map[string]int {
	namesOrdering := make(map[string]int)
	for i, name := range names {
		namesOrdering[name] = i
	}
	return namesOrdering
}

type HTTPIgnoreGetter struct {
	baseURL          string
	defaultExtension string
}

type RetrievedContents struct {
	namedSource NamedSource
	contents    string
	err         error
}

type NamedSource struct {
	name   string
	source string
}

func (getter *HTTPIgnoreGetter) GetIgnoreFiles(names []string, contentsChannel chan RetrievedContents, requestsPending *sync.WaitGroup) {
	namedURLs := getter.NamesToUrls(names)
	for _, namedURL := range namedURLs {
		requestsPending.Add(1)
		log.Println("Retrieving", namedURL.source)
		go fetchIgnoreFile(namedURL, contentsChannel, requestsPending)
	}
}

func (getter *HTTPIgnoreGetter) NamesToUrls(names []string) []NamedSource {
	urls := make([]NamedSource, len(names))
	for i, name := range names {
		urls[i] = getter.nameToURL(name)
	}
	return urls
}

func (getter *HTTPIgnoreGetter) nameToURL(name string) NamedSource {
	nameWithExtension := getter.getNameWithExtension(name)
	url := getter.baseURL + "/" + nameWithExtension
	return NamedSource{name, url}
}

func (getter *HTTPIgnoreGetter) getNameWithExtension(name string) string {
	if filepath.Ext(name) == "" {
		name = name + getter.defaultExtension
	}
	return name
}

type FailedSource struct {
	source string
	err    error
}

func (fs *FailedSource) Error() string {
	return fmt.Sprintf("%s %s", fs.source, fs.err.Error())
}

func fetchIgnoreFile(namedURL NamedSource, contentsChannel chan RetrievedContents, requestsPending *sync.WaitGroup) {
	defer requestsPending.Done()
	var fc RetrievedContents
	url := namedURL.source
	response, err := http.Get(url)
	if err != nil {
		fc = RetrievedContents{namedURL, "", err}
	} else if response.StatusCode != 200 {
		fc = RetrievedContents{namedURL, "", fmt.Errorf("Got status code %d", response.StatusCode)}
	} else {
		defer response.Body.Close()
		content, err := getContent(response.Body)
		if err != nil {
			fc = RetrievedContents{namedURL, "", fmt.Errorf("Error reading response body: %s", err.Error())}
		} else {
			fc = RetrievedContents{namedURL, content, nil}
		}
	}
	contentsChannel <- fc
}

func getContent(body io.ReadCloser) (content string, err error) {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		content = content + fmt.Sprintln(scanner.Text())
	}
	err = scanner.Err()
	return content, err
}

type FailedSources struct {
	sources []*FailedSource
}

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

type NamedIgnoreContents struct {
	name     string
	contents string
}

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
			failedSource := &FailedSource{retrievedContents.namedSource.source, retrievedContents.err}
			failedSources.Add(failedSource)
		} else {
			name := retrievedContents.namedSource.name
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
	app.Version = "0.1.0"
	app.Usage = "Creates gitignore files from central sources"

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "get",
			Usage: "Fetches gitignore patterns files from a central source and concatenates them",
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
					Usage: "Path to output file",
					Value: ".gitignore",
				},
			},
			ArgsUsage: "[gitignore_name] [gitignore_name …]",
			Action:    fetchAllIgnoreFiles,
		},
	}

	return app
}

func fetchAllIgnoreFiles(context *cli.Context) error {
	names := getNamesFromArguments(context)
	namesOrdering := createNamesOrdering(names)
	getter := HTTPIgnoreGetter{context.String("base-url"), context.String("default-extension")}
	contentsChannel := make(chan RetrievedContents, context.Int("max-connections"))
	var requestsPending sync.WaitGroup
	getter.GetIgnoreFiles(names, contentsChannel, &requestsPending)
	requestsPending.Wait()
	close(contentsChannel)
	contents, err := processContents(contentsChannel, namesOrdering)
	if err != nil {
		return err
	}
	outputFilePath := context.String("o")
	f, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	log.Println("Writing contents to", outputFilePath)
	err = writeIgnoreFile(f, contents)
	if err != nil {
		return err
	}
	log.Print("Finished")
	return err
}

func main() {
	log.SetFlags(0)
	app := creatCLI()
	app.RunAndExitOnError()
}
