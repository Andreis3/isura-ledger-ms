package uow

import (
	"context"
	"errors"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBeginTransaction    = errors.New("error opening transaction")
	ErrCommitTransaction   = errors.New("error committing transaction")
	ErrRollbackTransaction = errors.New("error rolling back transaction")
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return fault.BeginTransactionError(errors.Join(err, ErrBeginTransaction))
	}

	if err := fn(database.WithTx(ctx, tx)); err != nil {
		// Cria um contexto seguro para o rollback caso o ctx original tenha expirado/cancelado
		rbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if rbErr := tx.Rollback(rbCtx); rbErr != nil {
			return fault.RollbackTransactionError(errors.Join(err, rbErr, ErrRollbackTransaction))
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fault.CommitTransactionError(errors.Join(err, ErrCommitTransaction))
	}

	return nil
}
