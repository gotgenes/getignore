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
)

// Getter provides an implementation to retrieve gitignore patterns from
// files available over HTTP
type HTTPGetter struct {
	BaseURL          string
	DefaultExtension string
	MaxConnections   int
}

// GetIgnoreFiles retrieves gitignore patterns files via HTTP and sends their contents
// over a channel. It registers each request made with a WaitGroup instance, so the
// responses can be awaited.
func (getter *HTTPGetter) GetIgnoreFiles(names []string, contentsChannel chan contentstructs.RetrievedContents, requestsPending *sync.WaitGroup) {
	namesChannel := make(chan string)
	for i := 0; i < getter.MaxConnections; i++ {
		go getter.downloadIgnoreFile(namesChannel, contentsChannel, requestsPending)
	}
	for _, name := range names {
		requestsPending.Add(1)
		namesChannel <- name
	}
	close(namesChannel)
}

func (getter *HTTPGetter) downloadIgnoreFile(namesChannel chan string, contentsChannel chan contentstructs.RetrievedContents, requestsPending *sync.WaitGroup) {
	for name := range namesChannel {
		url := getter.nameToURL(name)
		log.Println("Retrieving", url)
		response, err := http.Get(url)
		contents, err := getter.processResponse(response, err)
		contentsChannel <- contentstructs.RetrievedContents{name, url, contents, err}
		requestsPending.Done()
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
