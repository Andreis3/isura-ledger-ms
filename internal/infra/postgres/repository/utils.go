package repository

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/jackc/pgx/v5/pgconn"
)

func resolveDB(ctx context.Context, db database.Querier) database.Querier {
	if tx, ok := database.ExtractTx(ctx); ok {
		return tx
	}
	return db
}

func isUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}
