package nats

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/domain/event"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js jetstream.JetStream
}

var _ event.Publisher = (*Publisher)(nil)

func NewJetStreamPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{
		js: js,
	}
}

func (p *Publisher) Publish(ctx context.Context, event event.Event) error {
	payload, err := event.Payload()
	if err != nil {
		return err
	}

	// Na API moderna, o contexto (ctx) entra como o primeiro argumento
	_, err = p.js.Publish(ctx, event.SubjectName(), payload)
	if err != nil {
		return err
	}

	return nil
}
