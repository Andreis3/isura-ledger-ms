package dto

import (
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/balance"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
)

type CreateBalanceInput struct {
	AccountID string `json:"account_id"`
	Currency  string `json:"currency"`
}

func (b *CreateBalanceInput) NewBalanceDomain() (*balance.Balance, error) {
	now := time.Now()
	initAmount := int64(0)

	return balance.NewBalanceBuilder().
		WithID().
		WithAccountID(b.AccountID).
		WithAmount(initAmount, money.Currency(b.Currency)).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()
}
