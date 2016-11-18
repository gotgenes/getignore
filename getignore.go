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

type HTTPIgnoreGetter struct {
	baseURL          string
	defaultExtension string
}

type NamedURL struct {
	name string
	url  string
}

func (getter *HTTPIgnoreGetter) NamesToUrls(names []string) []NamedURL {
	urls := make([]NamedURL, len(names))
	for i, name := range names {
		urls[i] = getter.nameToURL(name)
	}
	return urls
}

func (getter *HTTPIgnoreGetter) nameToURL(name string) NamedURL {
	nameWithExtension := getter.getNameWithExtension(name)
	url := getter.baseURL + "/" + nameWithExtension
	return NamedURL{name, url}
}

func (getter *HTTPIgnoreGetter) getNameWithExtension(name string) string {
	if filepath.Ext(name) == "" {
		name = name + getter.defaultExtension
	}
	return name
}

func addNamesToChannel(names []string, namesChannel chan string) {
	for _, v := range names {
		namesChannel <- v
	}
	close(namesChannel)
}

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

type FetchedContents struct {
	namedURL NamedURL
	contents string
	err      error
}

func getIgnoreFiles(contentsChannel chan FetchedContents, namedURLs []NamedURL) {
	var wg sync.WaitGroup
	for _, namedURL := range namedURLs {
		wg.Add(1)
		log.Println("Retrieving", namedURL.url)
		go fetchIgnoreFile(namedURL, contentsChannel, &wg)
	}
	wg.Wait()
	close(contentsChannel)
}

type FailedURL struct {
	url string
	err error
}

func (failedURL *FailedURL) Error() string {
	return fmt.Sprintf("%s %s", failedURL.url, failedURL.err.Error())
}

func fetchIgnoreFile(namedURL NamedURL, contentsChannel chan FetchedContents, wg *sync.WaitGroup) {
	defer wg.Done()
	var fc FetchedContents
	url := namedURL.url
	response, err := http.Get(url)
	if err != nil {
		fc = FetchedContents{namedURL, "", err}
	} else if response.StatusCode != 200 {
		fc = FetchedContents{namedURL, "", fmt.Errorf("Got status code %d", response.StatusCode)}
	} else {
		defer response.Body.Close()
		content, err := getContent(response.Body)
		if err != nil {
			fc = FetchedContents{namedURL, "", fmt.Errorf("Error reading response body: %s", err.Error())}
		} else {
			fc = FetchedContents{namedURL, content, nil}
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

type FailedURLs struct {
	urls []*FailedURL
}

func (failedURLs *FailedURLs) Add(failedURL *FailedURL) {
	failedURLs.urls = append(failedURLs.urls, failedURL)
}

func (failedURLs *FailedURLs) Error() string {
	urlErrors := make([]string, len(failedURLs.urls))
	for i, failedURL := range failedURLs.urls {
		urlErrors[i] = failedURL.Error()
	}
	stringOfErrors := strings.Join(urlErrors, "\n")
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

func processContents(contentsChannel chan FetchedContents, namesOrdering map[string]int) ([]NamedIgnoreContents, error) {
	retrievedContents := make([]NamedIgnoreContents, len(namesOrdering))
	var err error
	failedURLs := new(FailedURLs)
	for fetchedContents := range contentsChannel {
		if fetchedContents.err != nil {
			failedURL := &FailedURL{fetchedContents.namedURL.url, fetchedContents.err}
			failedURLs.Add(failedURL)
		} else {
			name := fetchedContents.namedURL.name
			position, present := namesOrdering[name]
			if !present {
				return retrievedContents, fmt.Errorf("Could not find name %s in ordering", name)
			}
			retrievedContents[position] = NamedIgnoreContents{name, fetchedContents.contents}
		}
	}
	if len(failedURLs.urls) > 0 {
		err = failedURLs
	}
	return retrievedContents, err
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
	getter := HTTPIgnoreGetter{context.String("base-url"), context.String("default-extension")}
	names := getNamesFromArguments(context)
	namesOrdering := createNamesOrdering(names)
	urls := getter.NamesToUrls(names)
	contentsChannel := make(chan FetchedContents, context.Int("max-connections"))
	go getIgnoreFiles(contentsChannel, urls)
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
