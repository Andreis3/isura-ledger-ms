package transaction

import (
	"errors"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	"github.com/andreis3/isura-ledger-ms/internal/domain/shared"
	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

var (
	ErrInvalidDirection    = errors.New("invalid direction")
	ErrNegativeAmountValue = errors.New("amount cannot be negative")
	ErrAmountEqualZero     = errors.New("amount cannot be zero")
)

type Direction string

const (
	Credit Direction = "CREDIT"
	Debit  Direction = "DEBIT"
)

func (d Direction) IsValid() bool {
	switch d {
	case Credit, Debit:
		return true
	}
	return false
}

type EntryBuilder struct {
	id             EntryID
	idempotencyKey string
	direction      Direction
	amount         money.Money
	accountID      AccountID
	transactionID  TransactionID
	createdAt      time.Time
	eval           validator.Evaluator
}

func NewEntryBuilder() *EntryBuilder {
	return &EntryBuilder{}
}

type Entry struct {
	ID             EntryID
	TransactionID  TransactionID
	AccountID      AccountID
	IdempotencyKey string
	Direction      Direction
	Amount         money.Money
	CreatedAt      time.Time
}

func (b *EntryBuilder) WithID(id string) *EntryBuilder {
	b.eval.CheckField(validator.NotBlank(id), "id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(id), "id", "is not uuid")
	b.id = EntryID(id)
	return b
}

func (b *EntryBuilder) WithTransactionID(transactionID string) *EntryBuilder {
	b.eval.CheckField(validator.NotBlank(transactionID), "transaction_id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(transactionID), "transaction_id", "is not uuid")
	b.transactionID = TransactionID(transactionID)
	return b
}

func (b *EntryBuilder) WithAccountID(accountID string) *EntryBuilder {
	b.eval.CheckField(validator.NotBlank(accountID), "account_id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(accountID), "account_id", "is not uuid")
	b.accountID = AccountID(accountID)
	return b
}

func (b *EntryBuilder) WithIdempotencyKey(key string) *EntryBuilder {
	b.eval.CheckField(validator.NotBlank(key), "idempotency_key", "cannot be blank")
	b.idempotencyKey = key
	return b
}

func (b *EntryBuilder) WithDirection(direction string) *EntryBuilder {
	d := Direction(direction)
	b.eval.CheckField(d.IsValid(), "direction", "invalid direction")
	b.direction = d
	return b
}

func (b *EntryBuilder) WithAmount(amount money.Money) *EntryBuilder {
	b.eval.CheckField(amount.IsZero(), "amount", "cannot be zero")
	b.eval.CheckField(amount.IsNegative(), "amount", "cannot be negative")
	b.eval.CheckField(amount.Currency().IsValid(), "amount", "invalid currency")
	b.amount = amount
	return b
}

func (b *EntryBuilder) WithCreatedAT(createdAt time.Time) *EntryBuilder {
	if !createdAt.IsZero() && createdAt.After(time.Now()) {
		b.eval.CheckField(false, "created_at", "cannot be in the future")
	}
	b.createdAt = createdAt
	return b
}

func (b *EntryBuilder) Build() (*Entry, error) {
	if len(b.eval) > 0 {
		return nil, fault.InvalidEntityError(errors.New("invalid account entity"), b.eval)
	}

	now := time.Now()

	return &Entry{
		ID:             b.id,
		TransactionID:  b.transactionID,
		AccountID:      b.accountID,
		IdempotencyKey: b.idempotencyKey,
		Direction:      b.direction,
		Amount:         b.amount,
		CreatedAt:      shared.CoalesceTime(b.createdAt, now),
	}, nil
}
