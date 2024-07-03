package main

import (
	"github.com/calmera/jetscript/lang"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func docsCommand() *cli.Command {
	return &cli.Command{
		Name:  "docs",
		Usage: "Generate the docs",
		Description: `Generate the docs based on the configuration of jetscript components.
The resulting docs will be in the form of json files where each file represents a component. This makes
it easier to consume the docs in a programmatic way later on.`,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "the directory to output the docs to",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			env := lang.NewEnvironment()

			outDir := c.Path("output")
			color.Green("Generating docs to %s", outDir)

			if err := env.DumpComponents(outDir); err != nil {
				color.Red("Failed to generate docs: %v", err)
				return nil
			}

			color.Green("Docs generated")
			return nil
		},
	}
}
