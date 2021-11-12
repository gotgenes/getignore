package cmd

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gotgenes/getignore/github"
	"github.com/gotgenes/getignore/writers"
	"github.com/urfave/cli/v2"
)

var Get = &cli.Command{
	Name:  "get",
	Usage: "retrieves gitignore patterns files from a central source, combines them, and outputs them",
	Flags: append(commonFlags, []cli.Flag{
		&cli.StringFlag{
			Name:    "output-file",
			Aliases: []string{"o"},
			Usage:   "Path to output file (default: STDOUT)",
		},
		&cli.StringFlag{
			Name:    "names-file",
			Aliases: []string{"n"},
			Usage:   "Path to file containing names of gitignore patterns files",
		},
		&cli.IntFlag{
			Name:    "max-requests",
			Aliases: []string{"m"},
			Usage:   "The number of maximum connections to open for HTTP requests",
			Value:   github.DefaultMaxRequests,
		},
	}...),
	ArgsUsage: "path [path …]",
	Action:    getFiles,
}

func getFiles(ctx *cli.Context) error {
	names := getNamesFromArguments(ctx)
	getter, err := newGithubGetter(ctx)
	if err != nil {
		return err
	}
	contents, err := getter.Get(ctx.Context, names)
	if err != nil {
		return err
	}
	outputFilePath, outputFile, err := getOutputFile(ctx)
	if err != nil {
		return err
	}
	log.Println("Writing contents to", outputFilePath)
	err = writers.WriteIgnoreFile(outputFile, contents)
	if err != nil {
		return err
	}
	return nil
}

func getNamesFromArguments(c *cli.Context) []string {
	names := c.Args().Slice()

	if c.String("names-file") != "" {
		namesFile, _ := os.Open(c.String("names-file"))
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

func getOutputFile(c *cli.Context) (string, io.Writer, error) {
	outputFilePath := c.String("o")
	var (
		outputFile io.Writer
		err        error
	)
	if outputFilePath == "" {
		outputFilePath = "STDOUT"
		outputFile = os.Stdout
	} else {
		outputFile, err = os.Create(outputFilePath)
	}
	return outputFilePath, outputFile, err
}
