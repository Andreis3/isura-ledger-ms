package balance

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
)

var (
	ErrBalanceNotFound = errors.New("balance not found")
)

type Repository interface {
	Save(ctx context.Context, balance *Balance) error
	Find(ctx context.Context, params criteria.BalanceCriteria) (*Balance, error)
}
