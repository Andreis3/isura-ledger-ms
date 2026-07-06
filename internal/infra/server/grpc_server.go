package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcTransport "github.com/andreis3/isura-ledger-ms/internal/transport/grpc"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/interceptor"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	deps       *BaseDeps
}

func NewGRPCServer(
	deps *BaseDeps,
) *GRPCServer {

	return &GRPCServer{
		deps: deps,
	}
}

func (s *GRPCServer) Start() {
	start := time.Now()

	// GRPC server com interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor(s.deps.Log.SlogJSON()),
			interceptor.MetricsInterceptor(s.deps.Prom),
			interceptor.TracingInterceptor(s.deps.Tracer),
		),
	)

	// registra todos os módulos
	registry := grpcTransport.NewServerRegistry(grpcServer, grpcTransport.NewLedgerModule(s.buildLedgerServer()))
	registry.RegisterAll()
	reflection.Register(grpcServer)

	s.deps.Log.InfoText("GRPC server started",
		slog.String("port", s.deps.Cfg.Servers.GRPC.Port),
		slog.String("startup_time", time.Since(start).String()),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.deps.Cfg.Servers.GRPC.Port))
	if err != nil {
		s.deps.Log.CriticalText("grpc server failed to listen",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := grpcServer.Serve(lis); err != nil {
		s.deps.Log.CriticalText("grpc server failed to serve",
			slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func (s *GRPCServer) buildLedgerServer() *grpcTransport.LedgerServer {

	composer := NewComposer(s.deps)

	accountRepo := composer.BuildAccountRepo()

	// use cases
	createAccount := command.NewCreateAccount(accountRepo, s.deps.Log, s.deps.Tracer)

	// handlers
	createAccountHandler := handler.NewCreateAccountHandler(createAccount, s.deps.Log, s.deps.Tracer)

	// server
	ledgerServer := grpcTransport.NewLedgerServer(createAccountHandler)

	// server
	return ledgerServer
}

func (s *GRPCServer) Stop() {
	s.grpcServer.GracefulStop()
}
