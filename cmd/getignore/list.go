package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

var List = &cli.Command{
	Name:   "list",
	Usage:  "lists available gitignore patterns files",
	Flags:  commonFlags,
	Action: listIgnoreFiles,
}

func listIgnoreFiles(ctx context.Context, cmd *cli.Command) error {
	getter, err := newGithubGetter(cmd)
	if err != nil {
		return err
	}
	ignoreFiles, err := getter.List(ctx)
	if err != nil {
		return err
	}
	outputString := strings.Join(ignoreFiles, "\n")
	_, err = fmt.Println(outputString)
	return err
}
