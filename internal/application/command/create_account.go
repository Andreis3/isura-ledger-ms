package command

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

var ErrAccountAlreadyExists = errors.New("account already exists")

type CreateAccountInput struct {
	OwnerID           string
	AccountExternalID string
	AccountNumber     string
	TaxID             string
	AccountingType    string
	Currency          string
}

type CreateAccount struct {
	accountRepository account.Repository
	log               application.Logger
	tracer            application.Tracer
}

func NewCreateAccount(
	accountRepository account.Repository,
	log application.Logger,
	tracer application.Tracer,
) *CreateAccount {
	return &CreateAccount{
		accountRepository: accountRepository,
		log:               log,
		tracer:            tracer,
	}
}

func (c *CreateAccount) Execute(ctx context.Context, input CreateAccountInput) (string, error) {
	ctx, span := c.tracer.Start(ctx, "CreateAccount.Execute")
	tracerID := span.SpanContext().TraceID()
	defer span.End()

	c.log.InfoJSON("CreateAccount received request",
		slog.String("trace_id", tracerID),
		slog.Any("input", MaskInput[CreateAccountInput](input)),
	)

	accountEntity, err := c.validate(input)
	if err != nil {
		span.RecordError(err)
		c.log.CriticalJSON("CreateAccount failed to validate",
			append([]any{
				slog.String("trace_id", tracerID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	existing, err := c.accountRepository.FindByAccountExternalID(ctx, accountEntity.AccountExternalID)
	if err != nil && !errors.Is(err, account.ErrAccountNotFound) {
		c.log.CriticalJSON("CreateAccount failed to find account by external ID",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	if existing != nil {
		domainErr := fault.Wrap(fault.CodeConflict, "account already exists", ErrAccountAlreadyExists)
		c.log.WarnJSON("CreateAccount account already exists",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(domainErr)...)...,
		)
		return string(existing.ID), nil
	}

	accountID := account.AccountID(uuid.New().String())

	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to create account entity",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	err = c.accountRepository.Save(ctx, accountEntity)
	if err != nil {
		c.log.CriticalJSON("CreateAccount failed to save account",
			append([]any{
				slog.String("trace_id", tracerID),
				slog.String("account_external_id", accountEntity.AccountExternalID)},
				fault.Attrs(err)...)...,
		)
		return "", err
	}

	c.log.InfoJSON("CreateAccount account created successfully",
		slog.String("trace_id", tracerID),
		slog.String("account_id", string(accountID)),
		slog.String("account_external_id", accountEntity.AccountExternalID),
	)

	return string(accountID), nil
}

func (c *CreateAccount) validate(input CreateAccountInput) (*account.Account, error) {
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
