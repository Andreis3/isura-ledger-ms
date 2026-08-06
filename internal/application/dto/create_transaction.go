package dto

type CreateTransactionInput struct {
	IdempotencyKey  *string `json:"idempotency_key"`
	DebitAccountID  *string `json:"debit_account_id"`
	CreditAccountID *string `json:"credit_account_id"`
	Amount          *int64  `json:"amount"`
	Currency        *string `json:"currency"`
}

type CreateTransactionOutput struct {
	TransactionID *string `json:"transaction_id"`
}
