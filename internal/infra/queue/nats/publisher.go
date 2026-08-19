package nats

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/domain/event"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	js nats.JetStreamContext
}

var _ event.Publisher = (*Publisher)(nil)

func NewJetStreamPublisher(js nats.JetStreamContext) *Publisher {
	return &Publisher{
		js: js,
	}
}

func (p *Publisher) Publish(ctx context.Context, event event.Event) error {
	payload, err := event.Payload()
	if err != nil {
		return err
	}

	_, err = p.js.Publish(event.SubjectName(), payload)
	if err != nil {
		return err
	}

	return nil
}
