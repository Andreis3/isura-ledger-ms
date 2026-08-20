package repository

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/domain/balance"
	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/model"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
	"github.com/jackc/pgx/v5"
)

type BalanceRepository struct {
	db database.Querier
}

func NewBalanceRepository(db database.Querier) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}

func (r *BalanceRepository) Save(ctx context.Context, balance *balance.Balance) error {
	db := resolveDB(ctx, r.db)

	query := `INSERT INTO
			balances(
				id,
				account_id,
				amount,
				currency,
				created_at,
				updated_at)
			values ($1, $2, $3, $4, $5, $6)		
		`

	balanceModel := model.ToBalanceModel(balance)

	_, err := db.Exec(ctx, query,
		balanceModel.ID,
		balanceModel.AccountID,
		balanceModel.Amount,
		balanceModel.Currency,
		balanceModel.CreatedAt,
		balanceModel.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			return fault.SaveBalanceAlreadyExistsError(err)
		}
		return fault.SaveModelError(balanceModel.TableName(), err)
	}

	return nil
}

func (r *BalanceRepository) Find(ctx context.Context, params criteria.BalanceCriteria) (*balance.Balance, error) {
	db := resolveDB(ctx, r.db)

	baseQuery := `
		SELECT
			id,
			account_id,
			amount,
			currency,
			created_at,
			updated_at
		FROM
		    balances
		WHERE 1 = 1
	`

	query, args := criteria.GetBalanceCriteria(baseQuery, params)
	var balanceModel model.Balance

	err := db.QueryRow(ctx, query, args...).Scan(
		&balanceModel.ID,
		&balanceModel.AccountID,
		&balanceModel.Amount,
		&balanceModel.Currency,
		&balanceModel.CreatedAt,
		&balanceModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Retorna erro padronizado de não encontrado
		}
		return nil, fault.FindAccountError(err)
	}

	balanceDomain, err := model.ToBalanceDomain(balanceModel)
	if err != nil {
		return nil, err
	}

	return balanceDomain, nil
}
