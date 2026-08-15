package model

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/entity"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction struct {
	ID             pgtype.Text
	IdempotencyKey pgtype.Text
	Status         pgtype.Text
	Operation      pgtype.Text
	Amount         pgtype.Int8
	Currency       pgtype.Text
	CreatedAt      pgtype.Timestamptz
	UpdatedAt      pgtype.Timestamptz
}

func ToTransactionModel(domain *transaction.Transaction) Transaction {
	return Transaction{
		ID: pgtype.Text{
			String: domain.ID.String(),
			Valid:  true,
		},
		IdempotencyKey: pgtype.Text{
			String: domain.IdempotencyKey,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: string(domain.Status),
			Valid:  true,
		},
		Operation: pgtype.Text{
			String: string(domain.Operation),
			Valid:  true,
		},
		Amount: pgtype.Int8{
			Int64: domain.Amount.Amount(),
			Valid: true,
		},
		Currency: pgtype.Text{
			String: string(domain.Amount.Currency()),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  domain.CreatedAt,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  domain.UpdatedAt,
			Valid: true,
		},
	}
}

func ToTransactionDomain(model Transaction, entries []*transaction.Entry) (*transaction.Transaction, error) {
	amount, err := money.NewMoney(model.Amount.Int64, money.Currency(model.Currency.String))
	if err != nil {
		return nil, err
	}

	id, err := entity.NewID(model.ID.String)
	if err != nil {
		return nil, err
	}

	return transaction.NewTransactionBuilder().
		WithID(id.String()).
		WithIdempotencyKey(model.IdempotencyKey.String).
		WithStatus(transaction.TransactionStatus(model.Status.String)).
		WithAmount(amount).
		WithOperation(transaction.Operation(model.Operation.String)).
		WithCreatedAt(model.CreatedAt.Time).
		WithUpdatedAt(model.UpdatedAt.Time).
		WithEntries(entries).
		Build()

}
