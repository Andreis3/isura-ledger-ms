package transaction

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
)

type Repository interface {
	Save(ctx context.Context, transaction *Transaction) error
	Find(ctx context.Context, params criteria.TransactionCriteria) (*Transaction, error)
	ExistsByIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error)
}
