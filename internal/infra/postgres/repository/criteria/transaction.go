package criteria

import "strings"

type TransactionCriteria struct {
	ID                   *string
	IdempotencyKey       *string
	Status               *string
	AccountID            *string
	Type                 *string
	HasForUpdateSkipLock bool
	WithEntries          bool
}

func GetTransactionCriteria(baseQuery string, params TransactionCriteria) (string, []any) {
	// Pre-allocates the slice with the maximum capacity of arguments (5 filters + slack)
	args := make([]any, 0, 6)
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

	if params.IdempotencyKey != nil {
		sb.WriteString(" AND idempotency_key = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.IdempotencyKey)
		argCount++
	}

	if params.Status != nil {
		sb.WriteString(" AND status = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.Status)
		argCount++
	}

	if params.AccountID != nil {
		sb.WriteString(" AND account_id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.AccountID)
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
