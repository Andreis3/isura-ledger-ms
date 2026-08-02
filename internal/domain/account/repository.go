package account

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Repository interface {
	Save(ctx context.Context, account *Account) error
	FindAccount(ctx context.Context, params criteria.AccountCriteria) (*Account, error)
}
