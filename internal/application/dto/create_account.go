package dto

type CreateAccountInput struct {
	AccountExternalID string `json:"account_external_id"`
	AccountNumber     string `json:"account_number"`
	TaxID             string `json:"tax_id" sensitive:"true"`
	Currency          string `json:"currency"`
}

type CreateAccountOutput struct {
	AccountID *string `json:"account_id"`
}
