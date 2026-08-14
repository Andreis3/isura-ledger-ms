package model

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/entity"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
	"github.com/jackc/pgx/v5/pgtype"
)

type Entry struct {
	ID            pgtype.Text
	AccountID     pgtype.Text
	TransactionID pgtype.Text
	Direction     pgtype.Text
	Amount        pgtype.Int8
	Currency      pgtype.Text
	CreatedAt     pgtype.Timestamptz
}

func ToEntryModel(domain *transaction.Entry) Entry {
	return Entry{
		ID: pgtype.Text{
			String: domain.ID.String(),
			Valid:  true,
		},
		Direction: pgtype.Text{
			String: string(domain.Direction),
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
		AccountID: pgtype.Text{
			String: domain.AccountID,
			Valid:  true,
		},
		TransactionID: pgtype.Text{
			String: domain.TransactionID,
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  domain.CreatedAt,
			Valid: true,
		},
	}
}

func ToEntryDomain(model Entry) (*transaction.Entry, error) {
	amount, err := money.NewMoney(model.Amount.Int64, money.Currency(model.Currency.String))
	if err != nil {
		return nil, err
	}

	id, err := entity.NewID(model.ID.String)
	if err != nil {
		return nil, err
	}

	return transaction.NewEntryBuilder().
		WithID(id.String()).
		WithTransactionID(model.TransactionID.String).
		WithAccountExternalID(model.AccountID.String).
		WithDirection(transaction.Direction(model.Direction.String)).
		WithAmount(amount).
		WithCreatedAt(model.CreatedAt.Time).
		Build()
}
