package criteria

import (
	"strings"
)

type AccountCriteria struct {
	ID                   *string
	AccountExternalID    *string
	TaxID                *string
	AccountNumber        *string
	Currency             *string
	Type                 *string
	HasForUpdateSkipLock bool
}

func GetAccountCriteria(baseQuery string, params AccountCriteria) (string, []any) {
	// Pre-allocates the slice with the maximum capacity of arguments (5 filters + slack)
	args := make([]any, 0, 7)
	argCount := 1

	// Estimates the approximate size of the query to avoid reallocations in the Builder
	var sb strings.Builder
	sb.Grow(len(baseQuery) + 128)
	sb.WriteString(baseQuery)

	if params.ID != nil {
		sb.WriteString(" AND id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.ID)
		argCount++
	}

	if params.AccountExternalID != nil {
		sb.WriteString(" AND account_external_id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.AccountExternalID)
		argCount++
	}

	if params.TaxID != nil {
		sb.WriteString(" AND tax_id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.TaxID)
		argCount++
	}

	if params.AccountNumber != nil {
		sb.WriteString(" AND account_number = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.AccountNumber)
		argCount++
	}

	if params.Currency != nil {
		sb.WriteString(" AND currency = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.Currency)
		argCount++
	}

	if params.Type != nil {
		sb.WriteString(" AND type = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.Type)
		argCount++
	}

	if params.HasForUpdateSkipLock {
		sb.WriteString(" FOR UPDATE SKIP LOCKED")
	}

	sb.WriteString(" LIMIT 1")

	return sb.String(), args
}

// Ultra-fast and zero-allocation helper function to convert argument indexes up to 9 to string
func argNumToString(n int) string {
	switch n {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	default:
		// Fallback in case it exceeds 9 parameters
		return string(rune('0' + n))
	}
}
