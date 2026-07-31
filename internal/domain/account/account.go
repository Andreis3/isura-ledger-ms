package account

import (
	"errors"
	"strings"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/tax"
	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

var (
	ErrInvalidAccountingType = errors.New("invalid accounting type")
	ErrEmptyExternalID       = errors.New("external id cannot be empty")
)

type AccountID string
type AccountType string

const (
	Asset     AccountType = "ASSET"
	Liability AccountType = "LIABILITY"
	Revenue   AccountType = "REVENUE"
	Expense   AccountType = "EXPENSE"
)

func (a AccountType) IsValid() bool {
	switch a {
	case Asset, Liability, Revenue, Expense:
		return true
	}
	return false
}

// AccountBuilder constrói Account passo a passo com validação
type AccountBuilder struct {
	id                AccountID
	ownerID           string
	accountExternalID string
	accountNumber     string
	accountType       AccountType
	taxID             string
	currency          money.Currency
	createdAt         time.Time
	updatedAt         time.Time
	eval              validator.Evaluator
}

// NewAccountBuilder inicia a construção
func NewAccountBuilder() *AccountBuilder {
	return &AccountBuilder{}
}

type Account struct {
	ID                AccountID
	OwnerID           string
	AccountExternalID string
	AccountNumber     string
	TaxID             string
	AccountType       AccountType
	Currency          money.Currency
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// WithID define o ID (obrigatório)
func (b *AccountBuilder) WithID(id string) *AccountBuilder {
	b.eval.CheckField(validator.NotBlank(id), "id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(id), "id", "is not uuid")
	b.id = AccountID(id)
	return b
}

// WithOwnerID define o OwnerID (obrigatório)
func (b *AccountBuilder) WithOwnerID(ownerID string) *AccountBuilder {
	b.eval.CheckField(validator.NotBlank(ownerID), "owner_id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(ownerID), "owner_id", "is not uuid")
	b.ownerID = ownerID
	return b
}

// WithAccountExternalID define o ID externo (obrigatório)
func (b *AccountBuilder) WithAccountExternalID(externalID string) *AccountBuilder {
	b.eval.CheckField(validator.NotBlank(externalID), "account_external_id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(externalID), "account_external_id", "is not uuid")
	b.accountExternalID = externalID
	return b
}

// WithAccountNumber define o número da conta (obrigatório)
func (b *AccountBuilder) WithAccountNumber(number string) *AccountBuilder {
	b.eval.CheckField(validator.NotBlank(number), "account_number", "cannot be blank")
	b.eval.CheckField(validator.MatchesNumber(number), "account_number", "is not a number")
	b.accountNumber = number
	return b
}

// WithAccountTaxID define CNPJ (obrigatório)
func (b *AccountBuilder) WithTaxID(rawTaxID string) *AccountBuilder {
	cnpjObj, eval := tax.NewCNPJ(rawTaxID)
	if len(eval) > 0 {
		// Mescla os erros do CNPJ no avaliador do Builder
		for field, msg := range eval {
			b.eval.AddFieldError(field, msg)
		}
		return b
	}
	b.taxID = cnpjObj.String()
	return b
}

// WithAccountType define o tipo da conta (obrigatório)
func (b *AccountBuilder) WithAccountType(accountType string) *AccountBuilder {
	accountTypeUpperCase := strings.ToUpper(accountType)
	b.eval.CheckField(validator.NotBlank(accountType), "account_type", "cannot be blank")

	at := AccountType(accountTypeUpperCase)
	b.eval.CheckField(at.IsValid(), "account_type", "invalid account type")
	b.accountType = at
	return b
}

// WithCurrency define a moeda (obrigatório)
func (b *AccountBuilder) WithCurrency(currency string) *AccountBuilder {
	currencyUpperCase := strings.ToUpper(currency)
	b.eval.CheckField(validator.NotBlank(currencyUpperCase), "currency", "cannot be blank")

	c := money.Currency(currencyUpperCase)
	b.eval.CheckField(c.IsValid(), "currency", "invalid currency")
	b.currency = c
	return b
}

// WithCreatedAt define a data de criação (opcional)
func (b *AccountBuilder) WithCreatedAt(createdAt time.Time) *AccountBuilder {
	if !createdAt.IsZero() && createdAt.After(time.Now()) {
		b.eval.CheckField(false, "created_at", "cannot be in the future")
	}
	b.createdAt = createdAt
	return b
}

// WithUpdatedAt define a data de atualização (opcional)
func (b *AccountBuilder) WithUpdatedAt(updatedAt time.Time) *AccountBuilder {
	if !updatedAt.IsZero() && updatedAt.After(time.Now()) {
		b.eval.CheckField(false, "updated_at", "cannot be in the future")
	}
	b.updatedAt = updatedAt
	return b
}

// Build constrói e valida a Account
func (b *AccountBuilder) Build() (*Account, error) {
	if len(b.eval) > 0 {
		return nil, fault.InvalidEntityError(errors.New("invalid account entity"), b.eval)
	}

	now := time.Now()

	return &Account{
		ID:                b.id,
		OwnerID:           b.ownerID,
		AccountExternalID: b.accountExternalID,
		AccountNumber:     b.accountNumber,
		AccountType:       b.accountType,
		TaxID:             b.taxID,
		Currency:          b.currency,
		CreatedAt:         coalesceTime(b.createdAt, now),
		UpdatedAt:         coalesceTime(b.updatedAt, now),
	}, nil
}

// coalesceTime retorna t se não for zero, caso contrário fallback
func coalesceTime(t, fallback time.Time) time.Time {
	if t.IsZero() {
		return fallback
	}
	return t
}
