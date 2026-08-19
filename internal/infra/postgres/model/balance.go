package model

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/balance"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/jackc/pgx/v5/pgtype"
)

type Balance struct {
	ID        pgtype.Text        `db:"id"`
	AccountID pgtype.Text        `db:"account_id"`
	Amount    pgtype.Int8        `db:"amount"`
	Currency  pgtype.Text        `db:"currency"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
}

func (m *Balance) TableName() string {
	return "balances"
}

func ToBalanceModel(entity *balance.Balance) Balance {
	return Balance{
		ID: pgtype.Text{
			String: entity.ID().String(),
			Valid:  true,
		},
		AccountID: pgtype.Text{
			String: entity.AccountID(),
			Valid:  true,
		},
		Amount: pgtype.Int8{
			Int64: entity.Amount().Amount(),
			Valid: true,
		},
		Currency: pgtype.Text{
			String: string(entity.Amount().Currency()),
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time:  entity.CreatedAT(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  entity.UpdatedAT(),
			Valid: true,
		},
	}
}

func ToBalanceDomain(model Balance) (*balance.Balance, error) {
	currency := money.Currency(model.Currency.String)

	return balance.NewBalanceBuilder().
		WithID(model.ID.String).
		WithAccountID(model.AccountID.String).
		WithAmount(model.Amount.Int64, currency).
		WithCreatedAt(model.CreatedAt.Time).
		WithUpdatedAt(model.UpdatedAt.Time).
		Build()
}
