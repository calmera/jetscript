package main

import (
	"context"
	"fmt"
	"github.com/calmera/jetscript/lang"
	"github.com/calmera/jetscript/utils"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func execCommand() *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "execute jetscript inline",
		Args:      true,
		ArgsUsage: "the path of the jetscript in the object store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "subject",
				Usage: "the subject to consume from",
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
			if c.NArg() == 0 {
				return fmt.Errorf("path is required")
			}

			if !c.IsSet("subject") {
				return fmt.Errorf("subject is required")
			}

			scriptPath := c.Args().First()
			subject := c.String("subject")
			contextName := c.String("context")
			bucket := c.String("bucket")

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

			ctx := context.Background()
			obj, err := js.ObjectStore(ctx, bucket)
			if err != nil {
				return fmt.Errorf("failed to connect to object store: %w", err)
			}

			script, err := obj.GetString(ctx, scriptPath)
			if err != nil {
				return fmt.Errorf("failed to get script: %w", err)
			}

			env := lang.NewEnvironment()
			exec, err := env.Parse(script)
			if err != nil {
				return err
			}

			_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
				result, err := lang.Process(exec, msg)
				if err != nil {
					log.Printf("failed to process message: %v", err)
					return
				}

				if result != nil {
					for _, m := range result {
						if err := nc.PublishMsg(m); err != nil {
							log.Printf("failed to publish message: %v", err)
						}
					}
				}

				// -- only if all the previous steps were successful we can ack the message
				_ = msg.Ack()
			})
			if err != nil {
				return fmt.Errorf("failed to subscribe to subject: %w", err)
			}

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			<-sigs
			fmt.Println("Received an interrupt, exiting...")

			return nil
		},
	}
}
