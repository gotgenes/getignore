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
	contentsChannel := make(chan contentstructs.NamedIgnoreContents)
	errorsChannel := make(chan errors.FailedSource)
	namedContentsChannel := make(chan []contentstructs.NamedIgnoreContents)
	failedSourcesChannel := make(chan errors.FailedSources)
	var jobsProcessing sync.WaitGroup
	jobsProcessing.Add(len(names))
	go getter.downloadIgnoreFiles(namesChannel, contentsChannel, errorsChannel)
	go processContents(contentsChannel, namedContentsChannel, &jobsProcessing)
	go processErrors(errorsChannel, failedSourcesChannel, &jobsProcessing)
	for _, name := range names {
		namesChannel <- name
	}
	close(namesChannel)
	jobsProcessing.Wait()
	close(contentsChannel)
	close(errorsChannel)
	contents = <-namedContentsChannel
	failedSources := <-failedSourcesChannel
	if len(failedSources) > 0 {
		err = failedSources
	}
	return
}

func (getter *HTTPGetter) downloadIgnoreFiles(namesChannel chan string, contentsChannel chan contentstructs.NamedIgnoreContents, failedSourceChannel chan errors.FailedSource) {
	for name := range namesChannel {
		go getter.downloadIgnoreFile(name, contentsChannel, failedSourceChannel)
	}
}

func (getter *HTTPGetter) downloadIgnoreFile(name string, contentsChannel chan contentstructs.NamedIgnoreContents, failedSourceChannel chan errors.FailedSource) {
	url := getter.nameToURL(name)
	log.Println("Retrieving", url)
	response, err := http.Get(url)
	contents, err := getter.processResponse(response, err)
	if err != nil {
		failedSourceChannel <- errors.FailedSource{url, err}
	} else {
		contentsChannel <- contentstructs.NamedIgnoreContents{name, contents}
	}
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

func processContents(contentsChannel chan contentstructs.NamedIgnoreContents, outputChannel chan []contentstructs.NamedIgnoreContents, jobsProcessing *sync.WaitGroup) {
	var allRetrievedContents []contentstructs.NamedIgnoreContents
	for retrievedContents := range contentsChannel {
		allRetrievedContents = append(allRetrievedContents, contentstructs.NamedIgnoreContents{retrievedContents.Name, retrievedContents.Contents})
		jobsProcessing.Done()
	}
	outputChannel <- allRetrievedContents
}

func processErrors(failedSourceChannel chan errors.FailedSource, collectedErrorsChannel chan errors.FailedSources, jobsProcessing *sync.WaitGroup) {
	var failedSources errors.FailedSources
	for failedSource := range failedSourceChannel {
		failedSources = append(failedSources, failedSource)
		jobsProcessing.Done()
	}
	collectedErrorsChannel <- failedSources
}
