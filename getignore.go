package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ignoreFetcher struct {
	baseUrl string
}

func (fetcher *ignoreFetcher) NamesToUrls(namesChannel chan string, urlsChannel chan string) {
	for name := range namesChannel {
		url := fetcher.NameToUrl(name)
		urlsChannel <- url
	}
	close(urlsChannel)
}

func (fetcher *ignoreFetcher) NameToUrl(name string) string {
	return fetcher.baseUrl + "/" + name + ".gitignore"
}

func addNamesToChannel(names []string, namesChannel chan string) {
	for _, v := range names {
		namesChannel <- v
	}
	close(namesChannel)
}

func getNamesFromArguments() []string {
	filePointer := flag.String("file", "", "Path to file of names")
	flag.Parse()
	names := flag.Args()

	if *filePointer != "" {
		namesFile, _ := os.Open(*filePointer)
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

func FetchIgnoreFiles(urlsChannel chan string, contentChannel chan string) error {
	var err error = nil
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
	var err error = nil
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

func main() {
	fetcher := ignoreFetcher{baseUrl: "https://raw.githubusercontent.com/github/gitignore/master"}
	names := getNamesFromArguments()
	namesChannel := make(chan string)
	urlsChannel := make(chan string)
	contentChannel := make(chan string)
	go addNamesToChannel(names, namesChannel)
	go fetcher.NamesToUrls(namesChannel, urlsChannel)
	go FetchIgnoreFiles(urlsChannel, contentChannel)
	f, _ := os.Create(".gitignore")
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	writeIgnoreFile(f, contentChannel, &waitGroup)
	waitGroup.Wait()
}
