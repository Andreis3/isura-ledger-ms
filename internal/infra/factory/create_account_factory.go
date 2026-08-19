package factory

import (
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/handler"
)

func NewCreateAccountFactory(
	baseDeps *dependency.BaseDeps,
) *handler.CreateAccountHandler {
	composeBuild := dependency.NewComposer(baseDeps)
	accountCommand := command.NewCreateAccount(
		composeBuild.BuildAccountRepo(),
		composeBuild.BuildNatsPublisher(),
		baseDeps.Log,
		baseDeps.Tracer,
		baseDeps.Prom,
	)

	createAccountHandler := handler.NewCreateAccountHandler(accountCommand, baseDeps.Log, baseDeps.Tracer)

	return createAccountHandler
}
