package observability

import (
	"context"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/domain/balance"
	"github.com/andreis3/isura-ledger-ms/internal/infra/postgres/repository/criteria"
)

type ObservabilityBalanceRepo struct {
	repo   balance.Repository
	metric application.Metrics
	tracer application.Tracer
}

func NewObservabilityBalanceRepo(repo balance.Repository, metric application.Metrics, tracer application.Tracer) *ObservabilityBalanceRepo {
	return &ObservabilityBalanceRepo{
		repo:   repo,
		metric: metric,
		tracer: tracer,
	}
}

func (r *ObservabilityBalanceRepo) Save(ctx context.Context, balance *balance.Balance) error {
	ctx, span := r.tracer.Start(ctx, "BalanceRepository.Save")
	defer span.End()

	start := time.Now()
	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"balances",
			"save",
			float64(time.Since(start).Milliseconds()))
	}()

	err := r.repo.Save(ctx, balance)

	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (r *ObservabilityBalanceRepo) Find(ctx context.Context, params criteria.BalanceCriteria) (*balance.Balance, error) {
	ctx, span := r.tracer.Start(ctx, "BalanceRepository.FindBalance")
	defer span.End()
	start := time.Now()

	defer func() {
		r.metric.RecordDBQueryDuration(
			"postgres",
			"balances",
			"find",
			float64(time.Since(start).Milliseconds()))
	}()

	balance, err := r.repo.Find(ctx, params)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return balance, nil
}
