package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var List = &cli.Command{
	Name:   "list",
	Usage:  "lists available gitignore patterns files",
	Flags:  commonFlags,
	Action: listIgnoreFiles,
}

func listIgnoreFiles(c *cli.Context) error {
	getter, err := newGithubGetter(c)
	if err != nil {
		return err
	}
	ctx := context.Background()
	ignoreFiles, err := getter.List(ctx)
	if err != nil {
		return err
	}
	outputString := strings.Join(ignoreFiles, "\n")
	_, err = fmt.Println(outputString)
	return err
}
