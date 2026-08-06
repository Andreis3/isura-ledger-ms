package transaction

import (
	"errors"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/validator"
)

var (
	ErrInvalidMaxEntries        = errors.New("maximum entries exceeded")
	ErrDuplicateEntryDirection  = errors.New("duplicate entry direction")
	ErrInvalidTransactionStatus = errors.New("invalid transaction status")
	ErrInvalidDifferentAmount   = errors.New("different amount")
	ErrTransactionNotFound      = errors.New("transaction not found")
)

type TransactionStatus string

const (
	Pending   TransactionStatus = "PENDING"
	Completed TransactionStatus = "COMPLETED"
	Failed    TransactionStatus = "FAILED"
)

func (t TransactionStatus) IsValid() bool {
	switch t {
	case Pending, Completed, Failed:
		return true
	}
	return false
}

type TransactionBuilder struct {
	id             TransactionID
	idempotencyKey string
	status         TransactionStatus
	entries        []*Entry
	createdAt      time.Time
	updatedAt      time.Time
	eval           validator.Evaluator
}

type Transaction struct {
	ID             TransactionID
	IdempotencyKey string
	Status         TransactionStatus
	Entries        []*Entry
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewTransactionBuilder initializes a new TransactionBuilder
func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{}
}

// WithID sets the transaction ID
func (b *TransactionBuilder) WithID(id string) *TransactionBuilder {
	b.eval.CheckField(validator.NotBlank(id), "id", "cannot be blank")
	b.eval.CheckField(validator.MatchesUUID(id), "id", "is not uuid")
	b.id = TransactionID(id)
	return b
}

// WithIdempotencyKey sets the idempotency key
func (b *TransactionBuilder) WithIdempotencyKey(key string) *TransactionBuilder {
	b.eval.CheckField(validator.NotBlank(key), "idempotency_key", "cannot be blank")
	b.idempotencyKey = key
	return b
}

// WithStatus sets the transaction status
func (b *TransactionBuilder) WithStatus(status string) *TransactionBuilder {
	s := TransactionStatus(status)
	b.eval.CheckField(s.IsValid(), "status", "invalid transaction status")
	b.status = s
	return b
}

// WithEntries adds entries to the transaction
func (b *TransactionBuilder) WithEntries(entries []*Entry) *TransactionBuilder {
	if len(entries) > 2 {
		b.eval.CheckField(false, "entries", "maximum entries exceeded")
		return b
	}

	if entries[0].Direction == entries[1].Direction {
		b.eval.CheckField(false, "entries", "duplicate entry direction: entries must have opposite directions (debit and credit)")
		return b
	}

	if entries[0].Amount != entries[1].Amount {
		b.eval.CheckField(false, "entries", "different amount")
		return b
	}

	b.entries = entries
	return b
}

func (t *Transaction) Complete() error {
	if !ValidStateMachine.CanTransition(t.Status, Completed) {
		return ErrInvalidTransactionStatus
	}

	t.Status = Completed
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) Fail() error {
	if !ValidStateMachine.CanTransition(t.Status, Failed) {
		return ErrInvalidTransactionStatus
	}

	t.Status = Failed
	t.UpdatedAt = time.Now()
	return nil
}
