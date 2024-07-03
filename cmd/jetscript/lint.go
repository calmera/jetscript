package main

import (
	"errors"
	"github.com/calmera/jetscript/lang"
	"github.com/fatih/color"
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/urfave/cli/v2"
	"os"
)

func lintCommand() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "lint a jetscript",
		Action: func(c *cli.Context) error {
			b, err := os.ReadFile(c.Args().First())
			if err != nil {
				return err
			}

			env := lang.NewEnvironment()
			if _, err = env.Parse(string(b)); err != nil {
				var pe *bloblang.ParseError
				if ok := errors.As(err, &pe); ok {
					color.Red(pe.ErrorMultiline())
					return nil
				}

				return err
			}

			color.Green("Jetscript is valid")
			return nil
		},
	}
}
