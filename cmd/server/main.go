package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
	"github.com/andreis3/isura-ledger-ms/internal/transport/queue/handler"

	_ "net/http/pprof" // <-- Automatically imports pprof handlers
)

func main() {
	deps := dependency.BuildBaseDeps()
	grpcSrv := server.NewGRPCServer(deps)
	httpSrv := server.NewHTTPServer(deps)
	// Cria o AccountConsumer usando a Factory (mesmo padrão das rotas HTTP)
	accountConsumer := handler.NewBalanceConsumer(deps)

	server.StartServersWithGracefulShutdown(grpcSrv, httpSrv, accountConsumer, deps)
}
