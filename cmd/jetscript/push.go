package main

import (
	"context"
	"fmt"
	"github.com/calmera/jetscript/utils"
	"github.com/fatih/color"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/urfave/cli/v2"
	"os"
)

func pushCommand() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "push the given jetscript file to the object store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Usage: "the resulting path in the object store",
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

			if !c.IsSet("path") {
				return fmt.Errorf("path is required")
			}

			path := c.String("path")
			contextName := c.String("context")
			bucket := c.String("bucket")

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

			// read the file
			f, err := os.Open(c.Args().First())
			if err != nil {
				return fmt.Errorf("failed to read the local jetscript file")
			}

			om := jetstream.ObjectMeta{
				Name: path,
				Metadata: map[string]string{
					"Content-Type": "application/x-jetscript",
					"Origin":       c.Args().First(),
				},
			}

			oi, err := obj.Put(ctx, om, f)
			if err != nil {
				return fmt.Errorf("failed to put object: %w", err)
			}

			color.Green("%s pushed to %s (%d bytes)", c.Args().First(), path, oi.Size)

			return nil
		},
	}
}
