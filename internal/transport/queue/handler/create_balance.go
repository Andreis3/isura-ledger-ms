package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application/dto"
	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/infra/factory"
	"github.com/nats-io/nats.go/jetstream"
)

type AccountConsumer struct {
	baseDeps *dependency.BaseDeps
}

func NewBalanceConsumer(baseDeps *dependency.BaseDeps) *AccountConsumer {
	return &AccountConsumer{
		baseDeps: baseDeps,
	}
}

func (c *AccountConsumer) Start(ctx context.Context) error {
	consumer, err := c.baseDeps.Nats.JS.CreateOrUpdateConsumer(ctx, "LEDGER_EVENTS", jetstream.ConsumerConfig{
		Name:          "BalanceCreationProcessor",
		Durable:       "BalanceCreationProcessor",
		FilterSubject: "ledger.events.account.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    5,
	})
	if err != nil {
		return err
	}

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		var input dto.CreateBalanceInput
		if err := json.Unmarshal(msg.Data(), &input); err != nil {
			c.baseDeps.Log.ErrorJSON("Failed to unmarshal account created event", "error", err.Error())
			_ = msg.Term()
			return
		}

		// === IGUAL AO HTTP: Instancia o Command fresco a cada mensagem via Factory ===
		createBalanceCommand := factory.NewCreateBalanceFactory(c.baseDeps)

		err := createBalanceCommand.Execute(ctx, input)
		if err != nil {
			c.baseDeps.Log.ErrorJSON("Failed to process create balance from event", "error", err.Error(), "account_id", input.AccountID)

			metadata, metaErr := msg.Metadata()
			if metaErr == nil && metadata.NumDelivered >= 5 {
				_ = msg.Term()
				return
			}
			_ = msg.NakWithDelay(5 * time.Second)
			return
		}

		_ = msg.Ack()
		c.baseDeps.Log.InfoJSON("Balance successfully created from event", "account_id", input.AccountID)
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	c.baseDeps.Log.InfoText("Stopping NATS AccountConsumer worker...")
	cc.Stop()
	return nil
}
