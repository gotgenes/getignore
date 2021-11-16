package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/gotgenes/getignore/cmd"
	"github.com/gotgenes/getignore/pkg/getignore"
)

func main() {
	log.SetFlags(0)
	app := creatCLI()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func creatCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "getignore"
	app.Version = getignore.Version
	app.Usage = "Bootstraps gitignore files from central sources"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		cmd.List,
		cmd.Get,
	}
	return app
}
