package main

import (
	"context"
	"fmt"
	"github.com/calmera/jetscript/utils"
	"github.com/fatih/color"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/urfave/cli/v2"
	"os"
	path2 "path"
)

func pullCommand() *cli.Command {
	return &cli.Command{
		Name:      "pull",
		Usage:     "pull a jetscript from the object store.",
		Args:      true,
		ArgsUsage: "the path of the jetscript in the object store",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "overwrite",
				Aliases: []string{"o"},
				Usage:   "overwrite the local file if it exists",
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "the name of the local file to write to",
			},
			&cli.StringFlag{
				Name:        "context",
				Usage:       "the nats context to use. Defaults to the currently selected context",
				DefaultText: natscontext.SelectedContext(),
			},
			&cli.StringFlag{
				Name:        "bucket",
				Usage:       "the object store bucket to use",
				DefaultText: "JETSCRIPT",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("path is required")
			}

			path := c.Args().First()
			target := c.String("target")
			contextName := c.String("context")
			bucket := c.String("bucket")
			overwrite := c.Bool("overwrite")

			nc, js, err := utils.ConnectJetstream(contextName)
			if err != nil {
				return err
			}
			defer nc.Close()

			ctx := context.Background()
			obj, err := js.ObjectStore(ctx, bucket)
			if err != nil {
				return fmt.Errorf("failed to connect to object store: %w", err)
			}

			// get the resulting filename from the path if a target is not set
			if target == "" {
				target = path2.Base(path)
			}

			// check if the file exists. If it does and the overwrite option has not been provided, return an error
			if _, err := os.Stat(target); err == nil && !overwrite {
				return fmt.Errorf("file %s already exists. Use the --overwrite flag to overwrite", target)
			}

			// get the data from the object store
			data, err := obj.GetBytes(ctx, path)
			if err != nil {
				return fmt.Errorf("failed to get object: %w", err)
			}

			// write the data to the file
			f, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to modify file: %w", err)
			}
			defer f.Close()
			cnt, err := f.Write(data)
			if err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			color.Green("%s pulled from %s (%d bytes)", target, path, cnt)

			return nil
		},
	}
}
