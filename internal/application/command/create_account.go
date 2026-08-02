package command

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application/dto"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
	"github.com/google/uuid"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

type CreateAccount struct {
	accountRepository account.Repository
	log               application.Logger
	tracer            application.Tracer
	metrics           application.Metrics
}

func NewCreateAccount(
	accountRepository account.Repository,
	log application.Logger,
	tracer application.Tracer,
	metrics application.Metrics,
) *CreateAccount {
	return &CreateAccount{
		accountRepository: accountRepository,
		log:               log,
		tracer:            tracer,
		metrics:           metrics,
	}
}

func (c *CreateAccount) Execute(ctx context.Context, input dto.CreateAccountInput) (*dto.CreateAccountOutput, error) {
	start := time.Now()
	ctx, span := c.tracer.Start(ctx, "CreateAccount.Execute")
	tracerID := span.SpanContext().TraceID()
	defer span.End()
	defer c.metrics.RecordCommandDuration("CreateAccount", float64(time.Since(start).Milliseconds()))

	c.log.InfoJSON("CreateAccount received request",
		slog.String("trace_id", tracerID),
		slog.Any("input", MaskInput[dto.CreateAccountInput](input)),
	)

	accountEntity, err := c.validate(input)
	if err != nil {
		span.RecordError(err)
		c.log.CriticalJSON("CreateAccount failed to validate",
			append([]any{
				slog.String("trace_id", tracerID)},
				fault.Attrs(err)...)...,
		)
		c.metrics.RecordCommandTotal("CreateAccount", "failure")
		return nil, err
	}

	parmsCriteria := criteria.AccountCriteria{
		AccountExternalID: &accountEntity.AccountExternalID,
		Currency:          new(string(accountEntity.Currency)),
	}

	existing, err := c.accountRepository.FindAccount(ctx, parmsCriteria)
	if err != nil && !errors.Is(err, account.ErrAccountNotFound) {
		c.log.CriticalJSON("CreateAccount failed to find account by external ID",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		c.metrics.RecordCommandTotal("CreateAccount", "failure")
		return nil, err
	}

	if existing != nil {
		c.log.InfoJSON("CreateAccount account already exists",
			slog.String("trace_id", tracerID),
			slog.String("account_external_id", accountEntity.AccountExternalID),
		)
		c.metrics.RecordCommandTotal("CreateAccount", "exist")
		return &dto.CreateAccountOutput{
			AccountID: new(string(existing.ID)),
		}, nil
	}

	accountID := account.AccountID(uuid.New().String())

	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to create account entity",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		c.metrics.RecordCommandTotal("CreateAccount", "failure")
		return nil, err
	}

	err = c.accountRepository.Save(ctx, accountEntity)
	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to save account",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		c.metrics.RecordCommandTotal("CreateAccount", "failure")
		return nil, err
	}

	c.log.InfoJSON("CreateAccount account created successfully",
		slog.String("trace_id", tracerID),
		slog.String("account_id", string(accountID)),
		slog.String("account_external_id", accountEntity.AccountExternalID),
	)

	return &dto.CreateAccountOutput{
		AccountID: new(string(accountID)),
	}, nil
}

func (c *CreateAccount) validate(input dto.CreateAccountInput) (*account.Account, error) {
	id := uuid.New().String()
	now := time.Now()

	return account.NewAccountBuilder().
		WithID(id).
		WithOwnerID(input.OwnerID).
		WithAccountExternalID(input.AccountExternalID).
		WithAccountNumber(input.AccountNumber).
		WithTaxID(input.TaxID).
		WithAccountType(input.AccountingType).
		WithCurrency(input.Currency).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
}
