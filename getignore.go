package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli"
)

type ignoreFetcher struct {
	baseURL string
}

func (fetcher *ignoreFetcher) NamesToUrls(namesChannel chan string, urlsChannel chan string) {
	for name := range namesChannel {
		url := fetcher.NameToURL(name)
		urlsChannel <- url
	}
	close(urlsChannel)
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

func fetchIgnoreFiles(urlsChannel chan string, contentChannel chan string) error {
	var err error
	for url := range urlsChannel {
		response, err := http.Get(url)
		if err != nil {
			close(contentChannel)
			return err
		}
		content, err := getContent(response.Body)
		if err != nil {
			close(contentChannel)
			return err
		}
		contentChannel <- content
	}
	close(contentChannel)
	return err
}

func getContent(body io.ReadCloser) (string, error) {
	var err error
	output := ""
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		output = output + fmt.Sprintln(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return output, err
}

func writeIgnoreFile(ignoreFile io.Writer, contentChannel chan string, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()
	for content := range contentChannel {
		_, err := io.WriteString(ignoreFile, content)
		if err != nil {
			fmt.Fprintln(os.Stderr, "writing output:", err)
			return err
		}
	}
	return nil
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
	namesChannel := make(chan string)
	urlsChannel := make(chan string)
	contentChannel := make(chan string)
	go addNamesToChannel(names, namesChannel)
	go fetcher.NamesToUrls(namesChannel, urlsChannel)
	go fetchIgnoreFiles(urlsChannel, contentChannel)
	f, _ := os.Create(context.String("o"))
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	writeIgnoreFile(f, contentChannel, &waitGroup)
	waitGroup.Wait()
	return nil
}

func main() {
	app := creatCLI()
	app.Run(os.Args)
}
