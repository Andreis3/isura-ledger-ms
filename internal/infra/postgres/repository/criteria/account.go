package criteria

import (
	"fmt"
)

type AccountCriteria struct {
	AccountExternalID *string
	TaxID             *string
	AccountNumber     *string
	AccountType       *string
	Currency          *string
}

func GetCriteria(query string, params AccountCriteria) (string, []any) {
	var args []any
	argCount := 1

	if params.AccountExternalID != nil {
		query += fmt.Sprintf(" AND account_external_id = $%d", argCount)
		args = append(args, *params.AccountExternalID)
		argCount++
	}

	if params.TaxID != nil {
		query += fmt.Sprintf(" AND tax_id = $%d", argCount)
		args = append(args, *params.TaxID)
		argCount++
	}

	if params.AccountNumber != nil {
		query += fmt.Sprintf(" AND account_number = $%d", argCount)
		args = append(args, *params.AccountNumber)
		argCount++
	}

	if params.AccountType != nil {
		query += fmt.Sprintf(" AND account_type = $%d", argCount)
		args = append(args, *params.AccountType)
		argCount++
	}

	if params.Currency != nil {
		query += fmt.Sprintf(" AND currency = $%d", argCount)
		args = append(args, *params.Currency)
		argCount++
	}

	query += " LIMIT 1"

	return query, args
}
