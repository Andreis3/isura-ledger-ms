package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
)

func main() {
	deps := server.BuildBaseDeps()
	grpcSrv := server.NewGRPCServer(deps)
	httpSrv := server.NewHTTPServer(*deps)

	server.StartServersWithGracefulShutdown(grpcSrv, httpSrv, *deps)
}
