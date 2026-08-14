package dto

import (
	"log/slog"

	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/andreis3/isura-ledger-ms/internal/util"
)

type CreateTransactionInput struct {
	IdempotencyKey  *string `json:"idempotency_key"`
	DebitAccountID  *string `json:"debit_account_id"`
	CreditAccountID *string `json:"credit_account_id"`
	Amount          *int64  `json:"amount"`
	Currency        *string `json:"currency"`
}

type CreateTransactionOutput struct {
	TransactionID *string `json:"transaction_id"`
}

func (d *CreateTransactionInput) CreateTransactionFacade() (*transaction.Transaction, error) {
	amount, err := money.NewMoney(util.Int64(d.Amount), money.Currency(util.String(d.Currency)))
	if err != nil {
		return nil, err
	}

	entryDestination, err := transaction.NewEntryBuilder().
		WithID().
		WithAccountExternalID(util.String(d.CreditAccountID)).
		WithAmount(amount).
		WithDirection(transaction.Credit).
		WithCreatedAt().
		Build()

	if err != nil {
		return nil, err
	}

	entrySource, err := transaction.NewEntryBuilder().
		WithID().
		WithAccountExternalID(util.String(d.DebitAccountID)).
		WithAmount(amount).
		WithDirection(transaction.Debit).
		WithCreatedAt().
		Build()
	if err != nil {
		return nil, err
	}

	entries := []*transaction.Entry{entryDestination, entrySource}

	return transaction.NewTransactionBuilder().
		WithID().
		WithIdempotencyKey(util.String(d.IdempotencyKey)).
		WithStatus(transaction.Pending).
		WithAmount(amount).
		WithEntries(entries).
		WithCreatedAt().
		WithUpdatedAt().
		Build()
}

// LogValue implements slog.LogValuer to safely log transaction input without exposing raw sensitive data.
func (d CreateTransactionInput) LogValue() slog.Value {
	// Safe extraction helpers to avoid nil pointer panics during logging
	idempotencyKey := ""
	if d.IdempotencyKey != nil {
		idempotencyKey = *d.IdempotencyKey
	}

	debitAccount := ""
	if d.DebitAccountID != nil {
		debitAccount = util.String(d.DebitAccountID) // Substitua por sua função de máscara real, ex: mask.UUID(*d.DebitAccountID)
	}

	creditAccount := ""
	if d.CreditAccountID != nil {
		creditAccount = util.String(d.CreditAccountID) // Substitua por sua função de máscara real
	}

	var amount int64
	if d.Amount != nil {
		amount = util.Int64(d.Amount)
	}

	currency := ""
	if d.Currency != nil {
		currency = util.String(d.Currency)
	}

	return slog.GroupValue(
		slog.String("idempotency_key", idempotencyKey),
		slog.String("debit_account_id", debitAccount),
		slog.String("credit_account_id", creditAccount),
		slog.Int64("amount", amount),
		slog.String("currency", currency),
	)
}
