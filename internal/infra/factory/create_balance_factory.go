package factory

import (
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
)

func NewCreateBalanceFactory(
	baseDeps *dependency.BaseDeps,
) *command.CreateBalance {
	composeBuild := dependency.NewComposer(baseDeps)
	balanceCommand := command.NewCreateBalance(
		composeBuild.BuildBalance(),
		composeBuild.BuildAccountRepo(),
		baseDeps.Log,
		baseDeps.Tracer,
		baseDeps.Prom,
	)

	return balanceCommand
}
