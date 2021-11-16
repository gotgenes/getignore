package main

import (
	"github.com/gotgenes/getignore/pkg/github"
	"github.com/urfave/cli/v2"
)

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "base-url",
		Aliases: []string{"u"},
		Usage:   "The base URL for the GitHub REST API v3 compatible server",
	},
	&cli.StringFlag{
		Name:    "owner",
		Aliases: []string{"w"},
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
}

var stringFlagsToOptions = map[string]func(string) github.GetterOption{
	"base-url":   github.WithBaseURL,
	"owner":      github.WithOwner,
	"repository": github.WithRepository,
	"branch":     github.WithBranch,
	"suffix":     github.WithSuffix,
}

func newGithubGetter(c *cli.Context) (github.Getter, error) {
	var opts []github.GetterOption
	for _, flagName := range c.FlagNames() {
		if flagName == "max-requests" {
			opts = append(opts, github.WithMaxRequests(c.Int(flagName)))
		} else {
			value := c.String(flagName)
			optFunc, ok := stringFlagsToOptions[flagName]
			if ok {
				opts = append(opts, optFunc(value))
			}
		}
	}
	getter, err := github.NewGetter(opts...)
	return getter, err
}
