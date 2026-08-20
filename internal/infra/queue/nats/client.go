package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type ClientNats struct {
	JS jetstream.JetStream
}

func NewJetStreamConnection(natsURL string) (*ClientNats, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	ctx := context.Background()

	if err := SetupStreams(ctx, js); err != nil {
		return nil, err
	}

	return &ClientNats{JS: js}, nil
}

func SetupStreams(ctx context.Context, js jetstream.JetStream) error {
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      "LEDGER_EVENTS",       // Nome padronizado com o Consumer
		Subjects:  []string{"ledger.>"},  // Captura qualquer evento que comece com ledger.
		Storage:   jetstream.FileStorage, // Persiste em disco
		Retention: jetstream.LimitsPolicy,
		MaxAge:    7 * 24 * time.Hour,
		Replicas:  1, // 1 réplica para ambiente local/Docker (evita erro de cluster)
		Discard:   jetstream.DiscardOld,
	})
	if err != nil {
		return fmt.Errorf("failed to create or update stream: %w", err)
	}
	return nil
}
