package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotgenes/getignore/github"
	"github.com/urfave/cli/v2"
)

var List = &cli.Command{
	Name:  "list",
	Usage: "Retrieves and prints a list of available ignore files",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "base-url",
			Aliases: []string{"u"},
			Usage:   "The base URL for the GitHub REST API v3 compatible server",
		},
		&cli.StringFlag{
			Name:    "owner",
			Aliases: []string{"o"},
			Usage:   "Owner/organization name of the gitignore repository",
			Value:   github.Owner,
		},
		&cli.StringFlag{
			Name:    "repository",
			Aliases: []string{"r"},
			Usage:   "Repository name of the gitignore repository",
			Value:   github.Repository,
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch or commit to inspect for the gitignore repository",
			Value:   github.Branch,
		},
		&cli.StringFlag{
			Name:    "suffix",
			Aliases: []string{"s"},
			Usage:   "The suffix to use to identify ignore files",
			Value:   github.Suffix,
		},
	},
	ArgsUsage: "",
	Action:    listIgnoreFiles,
}

var flagsToOptions = map[string]func(string) github.GitHubListerOption{
	"base-url":   github.WithBaseURL,
	"owner":      github.WithOwner,
	"repository": github.WithRepository,
	"branch":     github.WithBranch,
	"suffix":     github.WithSuffix,
}

func listIgnoreFiles(c *cli.Context) error {
	var opts []github.GitHubListerOption
	for flag, optFunc := range flagsToOptions {
		value := c.String(flag)
		if value != "" {
			opts = append(opts, optFunc(value))
		}
	}
	lister, err := github.NewGetter(opts...)
	if err != nil {
		return err
	}
	ctx := context.Background()
	ignoreFiles, err := lister.List(ctx)
	if err != nil {
		return err
	}
	outputString := strings.Join(ignoreFiles, "\n")
	_, err = fmt.Println(outputString)
	return err
}
