package repository

import (
	"context"
	"errors"

	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/database"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/model"
	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	db database.Querier
}

func NewAccountRepository(db database.Querier) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Save(ctx context.Context, account *account.Account) error {
	db := resolveDB(ctx, r.db)

	query := `INSERT INTO 
    accounts (id, 
              owner_id,
              account_external_id,
              account_number,
              tax_id,
              account_type,
              currency,
              created_at,
              updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	accountModel := model.ToAccountModel(account)

	_, err := db.Exec(ctx, query,
		accountModel.ID,
		accountModel.OwnerID,
		accountModel.AccountExternalID,
		accountModel.AccountNumber,
		accountModel.TaxID,
		accountModel.AccountType,
		accountModel.Currency,
		accountModel.CreatedAt,
		accountModel.UpdatedAt,
	)

	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT 
				a.id, 
				a.owner_id,
				a.account_external_id,
				a.account_number,
				a.tax_id,
				a.account_type,
				a.currency,
				a.created_at,
				a.updated_at
			FROM accounts a WHERE a.id = $1`

	var accountModel model.Account

	err := db.QueryRow(ctx, query, id).Scan(
		&accountModel.ID,
		&accountModel.OwnerID,
		&accountModel.AccountExternalID,
		&accountModel.AccountNumber,
		&accountModel.TaxID,
		&accountModel.AccountType,
		&accountModel.Currency,
		&accountModel.CreatedAt,
		&accountModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, account.ErrAccountNotFound
		}
		return nil, err
	}

	account, err := model.ToAccountDomain(accountModel)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountRepository) FindByAccountExternalID(ctx context.Context, AccountExternalID string) (*account.Account, error) {
	db := resolveDB(ctx, r.db)

	query := `SELECT 
				a.id, 
				a.owner_id,
				a.account_external_id,
				a.account_number,
				a.tax_id,
				a.account_type,
				a.currency,
				a.created_at,
				a.updated_at
			FROM accounts a WHERE a.account_external_id = $1`

	var accountModel model.Account

	err := db.QueryRow(ctx, query, AccountExternalID).Scan(
		&accountModel.ID,
		&accountModel.OwnerID,
		&accountModel.AccountExternalID,
		&accountModel.AccountNumber,
		&accountModel.TaxID,
		&accountModel.AccountType,
		&accountModel.Currency,
		&accountModel.CreatedAt,
		&accountModel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	account, err := model.ToAccountDomain(accountModel)
	if err != nil {
		return nil, err
	}

	return account, nil
}
