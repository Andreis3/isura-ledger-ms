package balance

import (
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/entity"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

type Balance struct {
	ID        entity.ID
	AccountID string
	Amount    money.Money
	CreatedAT time.Time
	UpdatedAt time.Time
}

type BalanceBuilder struct {
	id        entity.ID
	accountID string
	amount    money.Money
	createdAt time.Time
	updatedAt time.Time
	eval      validator.Evaluator
}

func NewBalanceBuilder() *BalanceBuilder {
	return &BalanceBuilder{}
}

func (b *BalanceBuilder) WithID(id ...string) *BalanceBuilder {
	if len(id) > 0 {
		b.eval.CheckField(validator.NotBlank(id[0]), "id", "cannot be blank")
		b.eval.CheckField(validator.MatchesUUID(id[0]), "id", "is not uuid")
		uuidV7, err := entity.NewID(id[0])
		if err != nil {
			b.eval.AddFieldError("id", err.Error())
		}

		b.id = uuidV7
		return b
	}

	uuidv7, err := entity.NewIDV7()
	if err != nil {
		b.eval.AddFieldError("id", err.Error())
	}
	b.id = uuidv7
	return b
}

func (b *BalanceBuilder) WithAccountID(accountID string) *BalanceBuilder {
	b.eval.CheckField(validator.NotBlank(accountID), "account_id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUIDv7(accountID), "account_id", "is not uuid")
	b.accountID = accountID
	return b
}

func (b *BalanceBuilder) WithAmount(amount int64, currency money.Currency) *BalanceBuilder {
	m, err := money.NewMoney(amount, currency)
	if err != nil {
		b.eval.AddFieldError("amount", err.Error())
	}
	b.amount = m
	return b
}

func (b *BalanceBuilder) WithCreatedAt(createdAt ...time.Time) *BalanceBuilder {
	if len(createdAt) > 0 {
		if !createdAt[0].IsZero() && createdAt[0].After(time.Now()) {
			b.eval.CheckField(false, "updated_at", "cannot be in the future")
		}
		b.createdAt = createdAt[0]
		return b
	}
	b.updatedAt = time.Now()
	return b
}

func (b *BalanceBuilder) WithUpdatedAt(updatedAt ...time.Time) *BalanceBuilder {
	if len(updatedAt) > 0 {
		if !updatedAt[0].IsZero() && updatedAt[0].After(time.Now()) {
			b.eval.CheckField(false, "updated_at", "cannot be in the future")
		}
		b.updatedAt = updatedAt[0]
		return b
	}
	b.updatedAt = time.Now()
	return b
}
