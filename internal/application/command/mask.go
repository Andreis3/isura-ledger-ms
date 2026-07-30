package command

import (
	"fmt"
	"log/slog"

	"github.com/andreis3/isura-ledger-ms/internal/application/dto"
)

func MaskInput[T any](input T) slog.Value {
	switch v := any(input).(type) {
	case dto.CreateAccountInput:
		return slog.GroupValue(
			slog.String("account_external_id", v.AccountExternalID),
			slog.String("accounting_type", v.AccountingType),
			slog.String("currency", v.Currency),
		)
	case CreateTransactionInput:
		return slog.GroupValue(
			slog.String("idempotency_key", v.IdempotencyKey),
			slog.String("debit_account_id", string(v.DebitAccountID)),
			slog.String("credit_account_id", string(v.CreditAccountID)),
			slog.String("amount", fmt.Sprintf("%d", v.Amount)),
			slog.String("currency", string(v.Currency)),
		)
	default:
		return slog.AnyValue(input)
	}
}
