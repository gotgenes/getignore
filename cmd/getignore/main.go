package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/gotgenes/getignore/pkg/getignore"
)

func main() {
	log.SetFlags(0)
	cmd := creatCLI()
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func creatCLI() *cli.Command {
	return &cli.Command{
		Name:                  "getignore",
		Version:               getignore.Version,
		Usage:                 "Bootstraps gitignore files from central sources",
		EnableShellCompletion: true,
		Commands:              []*cli.Command{List, Get},
	}
}
