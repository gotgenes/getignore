package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
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

func writeContent(writer bufio.Writer, contentChannel chan string) error {
	var err error = nil
	for content := range contentChannel {
		_, err := writer.WriteString(content)
		if err != nil {
			fmt.Fprintln(os.Stderr, "writing output:", err)
			return err
		}
	}
	return err
}

func main() {
}
