package handler

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
)

type CreateTransactionHandler struct {
	useCase *command.CreateTransaction
	log     application.Logger
	tracer  application.Tracer
}

func NewCreateTransactionHandler(
	useCase *command.CreateTransaction,
	log application.Logger,
	tracer application.Tracer,
) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		useCase: useCase,
		log:     log,
		tracer:  tracer,
	}
}

func (h *CreateTransactionHandler) Handle(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	ctx, span := h.tracer.Start(ctx, "CreateTransactionHandler.Handle")
	defer span.End()

	return nil, nil
}
