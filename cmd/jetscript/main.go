package main

import (
	_ "github.com/calmera/jetscript/components/nats"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := initApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initApp() *cli.App {
	app := &cli.App{
		Name:        "jetscript",
		Usage:       "Jetscript Commandline Tool",
		Description: "Jetscript is a tool for working with Jetscript files in combination with NATS Jetstream",
		Commands: []*cli.Command{
			lintCommand(),
			execCommand(),
			pushCommand(),
			pullCommand(),
		}}

	return app
}
