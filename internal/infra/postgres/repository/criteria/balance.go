package criteria

import "strings"

type BalanceCriteria struct {
	ID                   *string
	AccountID            *string
	HasForUpdateSkipLock bool
	WithEntries          bool
}

func GetBalanceCriteria(baseQuery string, params BalanceCriteria) (string, []any) {
	// Pre-allocates the slice with the maximum capacity of arguments (5 filters + slack)
	args := make([]any, 0, 3)
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

	if params.AccountID != nil {
		sb.WriteString(" AND account_id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.AccountID)
		argCount++
	}

	if params.HasForUpdateSkipLock {
		sb.WriteString(" FOR UPDATE SKIP LOCKED")
	}

	sb.WriteString(" LIMIT 1")

	return sb.String(), args
}
