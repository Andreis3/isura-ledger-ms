package event

import (
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application/dto"
	"github.com/andreis3/isura-ledger-ms/internal/domain/transaction"
)

type TransactionCreated struct {
	TransactionID   string    `json:"transaction_id"`
	IdempotencyKey  string    `json:"idempotency_key"`
	DebitAccountID  string    `json:"debit_account_id"`
	CreditAccountID string    `json:"credit_account_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	OccurredAt      time.Time `json:"occurred_at"`
}

func NewTransactionCreated() *TransactionCreated {
	return &TransactionCreated{}
}

func (t *TransactionCreated) WithTransactionID(transactionID string) *TransactionCreated {
	t.TransactionID = transactionID
	return t
}

func (t *TransactionCreated) WithIdempotencyKey(idempotencyKey string) *TransactionCreated {
	t.IdempotencyKey = idempotencyKey
	return t
}

func (t *TransactionCreated) WithDebitAccountID(debitAccountID string) *TransactionCreated {
	t.DebitAccountID = debitAccountID
	return t
}

func (t *TransactionCreated) WithCreditAccountID(creditAccountID string) *TransactionCreated {
	t.CreditAccountID = creditAccountID
	return t
}

func (t *TransactionCreated) WithAmount(amount int64) *TransactionCreated {
	t.Amount = amount
	return t
}

func (t *TransactionCreated) WithCurrency(currency string) *TransactionCreated {
	t.Currency = currency
	return t
}

func (t *TransactionCreated) WithStatus(status string) *TransactionCreated {
	t.Status = status
	return t
}

func (t *TransactionCreated) WithOccurredAt(occurredAt time.Time) *TransactionCreated {
	t.OccurredAt = occurredAt
	return t
}

func (t *TransactionCreated) Build() *TransactionCreated {
	return &TransactionCreated{
		TransactionID:   t.TransactionID,
		IdempotencyKey:  t.IdempotencyKey,
		DebitAccountID:  t.DebitAccountID,
		CreditAccountID: t.CreditAccountID,
		Amount:          t.Amount,
		Currency:        t.Currency,
		Status:          t.Status,
		OccurredAt:      t.OccurredAt,
	}
}

func TransactionCreatedFacade(entityTransaction transaction.Transaction, data dto.CreateTransactionInput) *TransactionCreated {
	return NewTransactionCreated().
		WithTransactionID(entityTransaction.ID.String()).
		WithIdempotencyKey(*data.IdempotencyKey).
		WithDebitAccountID(*data.DebitAccountID).
		WithCreditAccountID(*data.CreditAccountID).
		WithAmount(entityTransaction.Amount.Amount()).
		WithCurrency(string(entityTransaction.Amount.Currency())).
		WithStatus(string(entityTransaction.Status)).
		WithOccurredAt(time.Now()).
		Build()
}
