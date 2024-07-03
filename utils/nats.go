package utils

import (
	"fmt"
	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func ConnectJetstream(contextName string) (*nats.Conn, jetstream.JetStream, error) {
	nc, err := natscontext.Connect(contextName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to jetstream: %w", err)
	}

	return nc, js, nil
}
