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

type ignoreFetcher struct {
	baseURL string
}

func (fetcher *ignoreFetcher) NamesToUrls(names []string) []string {
	urls := make([]string, len(names))
	for i, name := range names {
		urls[i] = fetcher.NameToURL(name)
	}
	return urls
}

func (fetcher *ignoreFetcher) NameToURL(name string) string {
	return fetcher.baseURL + "/" + name + ".gitignore"
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

var MaxConnections int = 8

type FetchedContents struct {
	url      string
	contents string
	err      error
}

type NamedIgnoreContents struct {
	name     string
	contents string
}

func fetchIgnoreFiles(contentsChannel chan FetchedContents, urls []string) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		log.Println("Retrieving", url)
		go fetchIgnoreFile(url, contentsChannel, &wg)
	}
	wg.Wait()
	close(contentsChannel)
}

func fetchIgnoreFile(url string, contentsChannel chan FetchedContents, wg *sync.WaitGroup) {
	defer wg.Done()
	var fc FetchedContents
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fc = FetchedContents{url, "", fmt.Errorf("Error fetching URL %s", url)}
	} else {
		defer response.Body.Close()
		content, err := getContent(response.Body)
		if err != nil {
			fc = FetchedContents{url, "", fmt.Errorf("Error reading response body of %s", url)}
		} else {
			fc = FetchedContents{url, content, nil}
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

func processContents(contentsChannel chan FetchedContents) ([]NamedIgnoreContents, error) {
	var retrievedContents []NamedIgnoreContents
	var failedURLs []string
	var err error
	for fetchedContents := range contentsChannel {
		if fetchedContents.err != nil {
			failedURLs = append(failedURLs, fetchedContents.url)
		} else {
			retrievedContents = append(
				retrievedContents,
				NamedIgnoreContents{fetchedContents.url, fetchedContents.contents})
		}
	}
	if len(failedURLs) > 0 {
		err = fmt.Errorf("Failed to retrieve data from one or more URLs: %v", failedURLs)
	}
	return retrievedContents, err
}

func writeIgnoreFile(ignoreFile io.Writer, contents []NamedIgnoreContents) (err error) {
	for _, nc := range contents {
		_, err := io.WriteString(ignoreFile, nc.contents)
		if err != nil {
			fmt.Fprintln(os.Stderr, "writing output:", err)
			return err
		}
	}
	return
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
	fetcher := ignoreFetcher{baseURL: context.String("base-url")}
	names := getNamesFromArguments(context)
	urls := fetcher.NamesToUrls(names)
	contentsChannel := make(chan FetchedContents, MaxConnections)
	fetchIgnoreFiles(contentsChannel, urls)
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
