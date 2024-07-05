package main

import (
	"context"
	"github.com/calmera/jetscript/utils"
	"github.com/fatih/color"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/urfave/cli/v2"
)

func initCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "initialize the nats environment.",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "max-bytes",
				Usage:       "the maximum number of bytes that can be stored in the bucket",
				DefaultText: "-1",
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
			contextName := c.String("context")
			bucket := c.String("bucket")
			maxSize := c.Int64("max-bytes")

			if !c.IsSet("max-bytes") {
				maxSize = -1
			}

			if !c.IsSet("bucket") {
				bucket = "JETSCRIPT"
			}

			if !c.IsSet("context") {
				contextName = natscontext.SelectedContext()
			}

			nc, js, err := utils.ConnectJetstream(contextName)
			if err != nil {
				return err
			}
			defer nc.Close()

			cfg := jetstream.ObjectStoreConfig{
				Bucket:   bucket,
				MaxBytes: maxSize,
			}

			ctx := context.Background()
			if _, err := js.CreateOrUpdateObjectStore(ctx, cfg); err != nil {
				color.Red("failed to create JetStream %s Object Store: %s", bucket, err)
				return nil
			}

			color.Green("JetStream %s Object Store created", bucket)

			return nil
		},
	}
}
