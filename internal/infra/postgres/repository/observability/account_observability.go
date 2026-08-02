package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
)

type ObservabilityAccountRepo struct {
	repo   account.Repository
	metric application.Metrics
	tracer application.Tracer
}

func NewObservabilityAccountRepo(repo account.Repository, metric application.Metrics, tracer application.Tracer) *ObservabilityAccountRepo {
	return &ObservabilityAccountRepo{
		repo:   repo,
		metric: metric,
		tracer: tracer,
	}
}

func (r *ObservabilityAccountRepo) Save(ctx context.Context, account *account.Account) error {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.Save")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"save",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.Save(ctx, account)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityAccountRepo) FindAccount(ctx context.Context, params criteria.AccountCriteria) (*account.Account, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindByID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_id",
			float64(time.Since(start).Milliseconds()))
	}()

	account, err := r.repo.FindAccount(ctx, params)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return account, nil
}
