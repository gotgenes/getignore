package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/gotgenes/getignore/getters"
	"github.com/gotgenes/getignore/list"
	"github.com/gotgenes/getignore/writers"
)

// Version is the version of getignore
const Version string = "0.3.0.dev0"

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
	app.Version = Version
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
		cli.Command{
			Name:  "list",
			Usage: "Retrieves and prints a list of available ignore files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "api-url, u",
					Usage: "The GitHub Tree API-compatible URL to the repository of ignore files",
					Value: "https://api.github.com/repos/github/gitignore/git/trees/master?recursive=1",
				},
				cli.StringFlag{
					Name:  "suffix, s",
					Usage: "The suffix to use to identify ignore files",
					Value: ".gitignore",
				},
			},
			ArgsUsage: "",
			Action:    listIgnoreFiles,
		},
	}

	return app
}

func downloadAllIgnoreFiles(context *cli.Context) error {
	names := getNamesFromArguments(context)
	getter := getters.HTTPGetter{context.String("base-url"), context.String("default-extension"), context.Int("max-connections")}
	contents, err := getter.GetIgnoreFiles(names)
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

func listIgnoreFiles(context *cli.Context) error {
	outputString, err := list.ListIgnoreFiles(context.String("api-url"), Version, context.String("suffix"))
	if err != nil {
		return err
	}
	_, err = fmt.Println(outputString)
	return err
}

func main() {
	log.SetFlags(0)
	app := creatCLI()
	app.RunAndExitOnError()
}
