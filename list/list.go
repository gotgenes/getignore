package list

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const userAgentTemplate = "getignore/%s"
const acceptType string = "application/vnd.github.v3+json"

// ListIgnoreFiles downloads a list of ignore files and returns them as a
// newline-delimited string.  The function assumes the URL points to the GitHub
// Tree API, or an API compatible with the GitHub Tree API. If a suffix is
// provided, it is used to identify ignore files for display.
func ListIgnoreFiles(treeAPIURL string, applicationVersion string, suffix string) (namesString string, err error) {
	userAgentString := fmt.Sprintf(userAgentTemplate, applicationVersion)
	log.Println("Retrieving file names from", treeAPIURL)
	fileNames, err := getFileNamesFromTreeAPI(treeAPIURL, userAgentString)
	if err != nil {
		return
	}
	if suffix != "" {
		fileNames = filterBySuffix(fileNames, suffix)
	}
	namesString = strings.Join(fileNames, "\n")
	return
}

func getFileNamesFromTreeAPI(treeAPIURL string, userAgentString string) (fileNames []string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", treeAPIURL, nil)
	if err != nil {
		return
	}
	setRequestHeaders(request, userAgentString)
	response, err := client.Do(request)
	if err != nil {
		return
	} else if response.StatusCode != 200 {
		err = fmt.Errorf("Got status code %d", response.StatusCode)
		return
	}
	fileNames, err = parseGitTreeToFileNames(response.Body)
	if err != nil {
		return
	}
	err = response.Body.Close()
	return
}

func setRequestHeaders(request *http.Request, userAgentString string) {
	request.Header.Set("User-Agent", userAgentString)
	request.Header.Set("Accept", acceptType)
}

type treeInfo struct {
	Tree []fileInfo
}

type fileInfo struct {
	Path string
	Type string
}

func parseGitTreeToFileNames(reader io.Reader) (fileNames []string, err error) {
	var treeInfo treeInfo
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&treeInfo)
	for _, fileInfo := range treeInfo.Tree {
		if fileInfo.Type == "blob" {
			fileNames = append(fileNames, fileInfo.Path)
		}
	}
	return fileNames, err
}

func filterBySuffix(fileNames []string, suffix string) (filteredFileNames []string) {
	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, suffix) {
			filteredFileNames = append(filteredFileNames, fileName)
		}
	}
	return filteredFileNames
}
