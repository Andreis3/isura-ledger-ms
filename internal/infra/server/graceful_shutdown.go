package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(grpcSrv *GRPCServer, httpSrv *HTTPServer, deps BaseDeps) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	deps.Log.InfoText("isura-ledger-ms shutting down...")

	// 1. para de aceitar novos requests gRPC
	grpcSrv.Stop()
	deps.Log.InfoText("Stop grpc server...")

	// 2. para o HTTP com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Stop(ctx)
	deps.Log.InfoText("Stop http server...")

	// 3. fecha infraestrutura — só depois que os servers pararam
	deps.Pg.Close()
	deps.Log.InfoText("Close connection postgres...")

	deps.Prom.Close()
	deps.Log.InfoText("Close connection prometheus...")

	deps.TracerShutdown(ctx)
	deps.Log.InfoText("Close connection tracer...")

	deps.Log.InfoText("shutdown complete!")
	os.Exit(0)
}
