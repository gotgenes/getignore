package getters

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/errors"
)

// HTTPGetter provides an implementation to retrieve gitignore patterns from
// files available over HTTP
type HTTPGetter struct {
	BaseURL          string
	DefaultExtension string
	MaxConnections   int
}

// GetIgnoreFiles retrieves gitignore patterns files via HTTP and sends their contents
// over a channel. It registers each request made with a WaitGroup instance, so the
// responses can be awaited.
func (getter *HTTPGetter) GetIgnoreFiles(names []string) (contents []contentstructs.NamedIgnoreContents, err error) {
	namesChannel := make(chan string, getter.MaxConnections)
	contentsChannel := make(chan contentstructs.RetrievedContents)
	namedContentsChannel := make(chan contentsAndError)
	var jobsProcessing sync.WaitGroup
	jobsProcessing.Add(len(names))
	go getter.downloadIgnoreFiles(namesChannel, contentsChannel)
	go processContents(contentsChannel, namedContentsChannel, &jobsProcessing)
	for _, name := range names {
		namesChannel <- name
	}
	close(namesChannel)
	jobsProcessing.Wait()
	close(contentsChannel)
	results := <-namedContentsChannel
	contents = results.Contents
	err = results.Err
	return
}

func (getter *HTTPGetter) downloadIgnoreFiles(namesChannel chan string, contentsChannel chan contentstructs.RetrievedContents) {
	for name := range namesChannel {
		go getter.downloadIgnoreFile(name, contentsChannel)
	}
}

func (getter *HTTPGetter) downloadIgnoreFile(name string, contentsChannel chan contentstructs.RetrievedContents) {
	url := getter.nameToURL(name)
	log.Println("Retrieving", url)
	response, err := http.Get(url)
	contents, err := getter.processResponse(response, err)
	contentsChannel <- contentstructs.RetrievedContents{name, url, contents, err}
}

func (getter *HTTPGetter) nameToURL(name string) string {
	nameWithExtension := getter.getNameWithExtension(name)
	url := getter.BaseURL + "/" + nameWithExtension
	return url
}

func (getter *HTTPGetter) getNameWithExtension(name string) string {
	if filepath.Ext(name) == "" {
		name = name + getter.DefaultExtension
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

func (getter *HTTPGetter) processResponse(response *http.Response, err error) (contents string, processedErr error) {
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

type contentsAndError struct {
	Contents []contentstructs.NamedIgnoreContents
	Err      error
}

func processContents(contentsChannel chan contentstructs.RetrievedContents, outputChannel chan contentsAndError, jobsProcessing *sync.WaitGroup) {
	var allRetrievedContents []contentstructs.NamedIgnoreContents
	var err error
	failedSources := new(errors.FailedSources)
	for retrievedContents := range contentsChannel {
		if retrievedContents.Err != nil {
			failedSource := &errors.FailedSource{retrievedContents.Source, retrievedContents.Err}
			failedSources.Add(failedSource)
		} else {
			allRetrievedContents = append(allRetrievedContents, contentstructs.NamedIgnoreContents{retrievedContents.Name, retrievedContents.Contents})
		}
		jobsProcessing.Done()
	}
	if len(failedSources.Sources) > 0 {
		err = failedSources
	}
	outputChannel <- contentsAndError{Contents: allRetrievedContents, Err: err}
}
