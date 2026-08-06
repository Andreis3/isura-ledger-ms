package account

import (
	"errors"
	"strings"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/shared"
	"github.com/andreis3/isura-ledger-ms/internal/domain/tax"
	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

var (
	ErrInvalidAccountingType = errors.New("invalid accounting type")
	ErrEmptyExternalID       = errors.New("external id cannot be empty")
)

type AccountID string

type Type string

const (
	Asset     Type = "ASSET"     // Ativo (ex: Caixa, Contas no Banco Central, Direitos a receber)
	Liability Type = "LIABILITY" // Passivo (ex: Saldo de contas correntes dos clientes)
	Equity    Type = "EQUITY"    // Patrimônio Líquido (ex: Capital social)
	Revenue   Type = "REVENUE"   // Receitas (ex: Taxas de Pix e boletos cobradas)
	Expense   Type = "EXPENSE"   // Despesas (ex: Custos operacionais)
)

func (t Type) IsValid() bool {
	switch t {
	case Asset, Liability, Equity, Revenue, Expense:
		return true
	}
	return false
}

// AccountBuilder constrói Account passo a passo com validação
type AccountBuilder struct {
	id                AccountID
	accountExternalID string
	accountNumber     string
	taxID             string
	accountType       Type
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
	AccountExternalID string
	AccountNumber     string
	TaxID             string
	AccountType       Type
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

func (b *AccountBuilder) WithType(accountType string) *AccountBuilder {
	accountType = strings.ToUpper(accountType)
	t := Type(accountType)
	b.eval.CheckField(t.IsValid(), "type", "invalid account type")
	b.accountType = t
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
		AccountExternalID: b.accountExternalID,
		AccountNumber:     b.accountNumber,
		TaxID:             b.taxID,
		AccountType:       b.accountType,
		Currency:          b.currency,
		CreatedAt:         shared.CoalesceTime(b.createdAt, now),
		UpdatedAt:         shared.CoalesceTime(b.updatedAt, now),
	}, nil
}
