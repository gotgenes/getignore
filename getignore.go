package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/errors"
	"github.com/gotgenes/getignore/getters"
	"github.com/gotgenes/getignore/writers"
)

func getNamesFromArguments(context *cli.Context) []string {
	names := context.Args()

	if context.String("names-file") != "" {
		namesFile, _ := os.Open(context.String("names-file"))
		names = append(names, ParseNamesFile(namesFile)...)
	}
	return names
}

// ParseNamesFile reads a file containing one name of a gitignore patterns file per line
func ParseNamesFile(namesFile io.Reader) []string {
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

// CreateNamesOrdering creates a mapping of each name to its respective input position
func CreateNamesOrdering(names []string) map[string]int {
	namesOrdering := make(map[string]int)
	for i, name := range names {
		namesOrdering[name] = i
	}
	return namesOrdering
}

func processContents(contentsChannel chan contentstructs.RetrievedContents, namesOrdering map[string]int) ([]contentstructs.NamedIgnoreContents, error) {
	allRetrievedContents := make([]contentstructs.NamedIgnoreContents, len(namesOrdering))
	var err error
	failedSources := new(errors.FailedSources)
	for retrievedContents := range contentsChannel {
		if retrievedContents.Err != nil {
			failedSource := &errors.FailedSource{retrievedContents.Source, retrievedContents.Err}
			failedSources.Add(failedSource)
		} else {
			name := retrievedContents.Name
			position, present := namesOrdering[name]
			if !present {
				return allRetrievedContents, fmt.Errorf("Could not find name %s in ordering", name)
			}
			allRetrievedContents[position] = contentstructs.NamedIgnoreContents{name, retrievedContents.Contents}
		}
	}
	if len(failedSources.Sources) > 0 {
		err = failedSources
	}
	return allRetrievedContents, err
}

func getOutputFile(context *cli.Context) (outputFilePath string, outputFile io.Writer, err error) {
	outputFilePath = context.String("o")
	if outputFilePath == "" {
		outputFilePath = "STDOUT"
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(outputFilePath)
	}
	return
}

func creatCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "getignore"
	app.Version = "0.2.0.dev"
	app.Usage = "Bootstraps gitignore files from central sources"

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "get",
			Usage: "Retrieves gitignore patterns files from a central source and concatenates them",
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
					Usage: "Path to output file (default: STDOUT)",
				},
			},
			ArgsUsage: "[gitignore_name] [gitignore_name …]",
			Action:    downloadAllIgnoreFiles,
		},
	}

	return app
}

func downloadAllIgnoreFiles(context *cli.Context) error {
	names := getNamesFromArguments(context)
	namesOrdering := CreateNamesOrdering(names)
	getter := getters.HTTPGetter{context.String("base-url"), context.String("default-extension"), context.Int("max-connections")}
	contentsChannel := make(chan contentstructs.RetrievedContents, context.Int("max-connections"))
	var requestsPending sync.WaitGroup
	getter.GetIgnoreFiles(names, contentsChannel, &requestsPending)
	requestsPending.Wait()
	close(contentsChannel)
	contents, err := processContents(contentsChannel, namesOrdering)
	if err != nil {
		return err
	}
	outputFilePath, outputFile, err := getOutputFile(context)
	if err != nil {
		return err
	}
	log.Println("Writing contents to", outputFilePath)
	err = writers.WriteIgnoreFile(outputFile, contents)
	if err != nil {
		return err
	}
	return err
}

func main() {
	log.SetFlags(0)
	app := creatCLI()
	app.RunAndExitOnError()
}
