package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	JS nats.JetStreamContext
}

func NewJetStreamConnection(natsURL string) (*Nats, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	// Já garante a criação do Stream assim que conecta
	if err := SetupStreams(js); err != nil {
		return nil, fmt.Errorf("failed to setup nats streams: %w", err)
	}

	return &Nats{JS: js}, nil
}

func SetupStreams(js nats.JetStreamContext) error {
	// Cria um Stream chamado "LEDGER_EVENTS" que escuta tudo que começa com "ledger.events.>"
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "LEDGER_EVENTS",
		Subjects: []string{"ledger.events.>"},
		Storage:  nats.FileStorage, // Salva em disco dentro do container
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}
	return nil
}
