package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

type IgnoreFetcher struct {
	baseURL string
}

type NamedURL struct {
	name string
	url  string
}

func (fetcher *IgnoreFetcher) NamesToUrls(names []string) []NamedURL {
	urls := make([]NamedURL, len(names))
	for i, name := range names {
		urls[i] = fetcher.NameToURL(name)
	}
	return urls
}

func (fetcher *IgnoreFetcher) NameToURL(name string) NamedURL {
	url := fetcher.baseURL + "/" + name + ".gitignore"
	return NamedURL{name, url}
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

type FetchedContents struct {
	namedURL NamedURL
	contents string
	err      error
}

type NamedIgnoreContents struct {
	name     string
	contents string
}

func fetchIgnoreFiles(contentsChannel chan FetchedContents, namedURLs []NamedURL) {
	var wg sync.WaitGroup
	for _, namedURL := range namedURLs {
		wg.Add(1)
		log.Println("Retrieving", namedURL.url)
		go fetchIgnoreFile(namedURL, contentsChannel, &wg)
	}
	wg.Wait()
	close(contentsChannel)
}

func fetchIgnoreFile(namedURL NamedURL, contentsChannel chan FetchedContents, wg *sync.WaitGroup) {
	defer wg.Done()
	var fc FetchedContents
	url := namedURL.url
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fc = FetchedContents{namedURL, "", fmt.Errorf("Error fetching URL %s", url)}
	} else {
		defer response.Body.Close()
		content, err := getContent(response.Body)
		if err != nil {
			fc = FetchedContents{namedURL, "", fmt.Errorf("Error reading response body of %s", url)}
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
	URLs []string
}

func (failedURLs *FailedURLs) Add(url string) {
	failedURLs.URLs = append(failedURLs.URLs, url)
}

func (failedURLs *FailedURLs) Error() string {
	stringOfURLs := strings.Join(failedURLs.URLs, "\n")
	return "Failed to retrieve or read content from the following URLs:\n" + stringOfURLs
}

func processContents(contentsChannel chan FetchedContents) ([]NamedIgnoreContents, error) {
	var retrievedContents []NamedIgnoreContents
	var err error
	failedURLs := new(FailedURLs)
	for fetchedContents := range contentsChannel {
		if fetchedContents.err != nil {
			failedURLs.Add(fetchedContents.namedURL.url)
		} else {
			retrievedContents = append(
				retrievedContents,
				NamedIgnoreContents{fetchedContents.namedURL.name, fetchedContents.contents})
		}
	}
	if len(failedURLs.URLs) > 0 {
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
		writer.WriteString(decorateName(nc.name))
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
					Name:  "base-url",
					Usage: "The URL under which gitignore files can be found",
					Value: "https://raw.githubusercontent.com/github/gitignore/master",
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
	fetcher := IgnoreFetcher{baseURL: context.String("base-url")}
	names := getNamesFromArguments(context)
	urls := fetcher.NamesToUrls(names)
	contentsChannel := make(chan FetchedContents, context.Int("max-connections"))
	go fetchIgnoreFiles(contentsChannel, urls)
	contents, err := processContents(contentsChannel)
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
