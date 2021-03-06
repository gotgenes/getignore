package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/gotgenes/getignore/getters"
	"github.com/gotgenes/getignore/list"
	"github.com/gotgenes/getignore/writers"
)

// Version is the version of getignore.
// It should be populated through ldflags, e.g.,
// -ldflags "-X main.Version=${VERSION}"
var Version string

func main() {
	log.SetFlags(0)
	app := creatCLI()
	app.RunAndExitOnError()
}

func creatCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "getignore"
	app.Version = Version
	app.Usage = "Bootstraps gitignore files from central sources"
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		{
			Name:  "get",
			Usage: "Retrieves gitignore patterns files from a central source and concatenates them",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "base-url",
					Aliases: []string{"u"},
					Usage:   "The URL under which gitignore files can be found",
					Value:   "https://raw.githubusercontent.com/github/gitignore/master",
				},
				&cli.StringFlag{
					Name:    "default-extension",
					Aliases: []string{"e"},
					Usage:   "The default file extension appended to names when retrieving them",
					Value:   ".gitignore",
				},
				&cli.IntFlag{
					Name:  "max-connections",
					Usage: "The number of maximum connections to open for HTTP requests",
					Value: 8,
				},
				&cli.StringFlag{
					Name:    "names-file",
					Aliases: []string{"n"},
					Usage:   "Path to file containing names of gitignore patterns files",
				},
				&cli.StringFlag{
					Name:    "output-file",
					Aliases: []string{"o"},
					Usage:   "Path to output file (default: STDOUT)",
				},
			},
			ArgsUsage: "[gitignore_name] [gitignore_name …]",
			Action:    downloadAllIgnoreFiles,
		},
		{
			Name:  "list",
			Usage: "Retrieves and prints a list of available ignore files",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "api-url",
					Aliases: []string{"u"},
					Usage:   "The GitHub Tree API-compatible URL to the repository of ignore files",
					Value:   "https://api.github.com/repos/github/gitignore/git/trees/master?recursive=1",
				},
				&cli.StringFlag{
					Name:    "suffix",
					Aliases: []string{"s"},
					Usage:   "The suffix to use to identify ignore files",
					Value:   ".gitignore",
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
	getter := getters.HTTPGetter{
		BaseURL:          context.String("base-url"),
		DefaultExtension: context.String("default-extension"),
		MaxConnections:   context.Int("max-connections"),
	}
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

func getNamesFromArguments(context *cli.Context) []string {
	names := context.Args().Slice()

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

func listIgnoreFiles(context *cli.Context) error {
	outputString, err := list.ListIgnoreFiles(context.String("api-url"), Version, context.String("suffix"))
	if err != nil {
		return err
	}
	_, err = fmt.Println(outputString)
	return err
}
