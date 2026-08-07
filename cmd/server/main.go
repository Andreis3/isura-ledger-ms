package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/dependency"
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"

	_ "net/http/pprof" // <-- Automatically imports pprof handlers
)

func main() {
	deps := dependency.BuildBaseDeps()
	grpcSrv := server.NewGRPCServer(deps)
	httpSrv := server.NewHTTPServer(deps)

	server.StartServersWithGracefulShutdown(grpcSrv, httpSrv, deps)
}
