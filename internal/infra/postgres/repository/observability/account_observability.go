package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/account"
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

func (r *ObservabilityAccountRepo) FindByID(ctx context.Context, id account.AccountID) (*account.Account, error) {
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

	account, err := r.repo.FindByID(ctx, id)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return account, nil
}

func (r *ObservabilityAccountRepo) FindByAccountExternalID(ctx context.Context, externalID string) (*account.Account, error) {
	ctx, span := r.tracer.Start(ctx, "AccountRepository.FindByExternalID")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"accounts",
			"find_by_external_id",
			float64(time.Since(start).Milliseconds()))
	}()

	account, err := r.repo.FindByAccountExternalID(ctx, externalID)

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return account, nil
}
