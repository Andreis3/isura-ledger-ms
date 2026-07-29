package account

import (
	"context"
	"errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type Repository interface {
	Save(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id AccountID) (*Account, error)
	FindByAccountExternalID(ctx context.Context, externalID string) (*Account, error)
}
