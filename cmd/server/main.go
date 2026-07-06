package main

import (
	"github.com/andreis3/isura-ledger-ms/internal/infra/server"
)

func main() {
	deps := server.BuildBaseDeps()

	grpcSrv := server.NewGRPCServer(deps)

	httpSrv := server.NewHTTPServer(*deps)

	go httpSrv.Start()
	go grpcSrv.Start()

	server.GracefulShutdown(grpcSrv, httpSrv, *deps)
}
