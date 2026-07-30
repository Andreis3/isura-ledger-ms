package dto

type CreateAccountInput struct {
	OwnerID           string `json:"owner_id"`
	AccountExternalID string `json:"account_external_id"`
	AccountNumber     string `json:"account_number"`
	TaxID             string `json:"tax_id"`
	AccountingType    string `json:"accounting_type"`
	Currency          string `json:"currency"`
}
