package criteria

import "strings"

type EntryCriteria struct {
	TransactionID *string
}

func GetEntryCriteria(baseQuery string, params EntryCriteria) (string, []any) {
	args := make([]any, 0, 1)
	argCount := 1
	var sb strings.Builder

	sb.Grow(len(baseQuery) + 128)
	sb.WriteString(baseQuery)

	if params.TransactionID != nil {
		sb.WriteString(" AND transaction_id = $")
		sb.WriteString(argNumToString(argCount))
		args = append(args, *params.TransactionID)
		argCount++
	}

	sb.WriteString(" LIMIT 2")

	return sb.String(), args
}
