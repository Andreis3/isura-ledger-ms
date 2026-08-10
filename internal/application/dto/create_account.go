package dto

import (
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/google/uuid"
)

type CreateAccountInput struct {
	AccountExternalID string `json:"account_external_id"`
	AccountNumber     string `json:"account_number"`
	TaxID             string `json:"tax_id" sensitive:"partial"`
	AccountType       string `json:"account_type"`
	Currency          string `json:"currency"`
}

type CreateAccountOutput struct {
	AccountID *string `json:"account_id"`
}

func (d *CreateAccountInput) CreateAccountFacade(id ...string) (*account.Account, error) {
	if len(id) == 0 {
		id[0] = uuid.New().String()
	}
	return account.NewAccountBuilder().
		WithID(id[0]).
		WithAccountExternalID(d.AccountExternalID).
		WithAccountNumber(d.AccountNumber).
		WithTaxID(d.TaxID).
		WithStatus().
		WithType(d.AccountType).
		WithCurrency(d.Currency).
		WithCreatedAt().
		WithUpdatedAt().
		Build()
}
