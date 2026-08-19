package account

import (
	"encoding/json"
)

type AccountCreatedEvent struct {
	AccountID string `json:"account_id"`
	Currency  string `json:"currency"`
}

func NewAccountCreatedEvent(accountID, currency string) *AccountCreatedEvent {
	return &AccountCreatedEvent{
		AccountID: accountID,
		Currency:  currency,
	}
}

// SubjectName define o canal/tópico no NATS JetStream
func (e *AccountCreatedEvent) SubjectName() string {
	return "ledger.events.account.created"
}

// Payload serializa a struct para JSON
func (e *AccountCreatedEvent) Payload() ([]byte, error) {
	return json.Marshal(e)
}
