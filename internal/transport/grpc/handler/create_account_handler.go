package handler

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
	"github.com/andreis3/isura-ledger-ms/internal/application/dto"
	"github.com/andreis3/isura-ledger-ms/internal/domain/money"
	pb "github.com/andreis3/isura-ledger-ms/internal/transport/grpc/pb/ledger/v1"
	"github.com/andreis3/isura-ledger-ms/internal/transport/grpc/translator"
)

type CreateAccountHandler struct {
	useCase *command.CreateAccount
	log     application.Logger
	tracer  application.Tracer
}

func NewCreateAccountHandler(
	useCase *command.CreateAccount,
	log application.Logger,
	tracer application.Tracer,
) *CreateAccountHandler {
	return &CreateAccountHandler{
		useCase: useCase,
		log:     log,
		tracer:  tracer,
	}
}

func (h *CreateAccountHandler) Handle(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	ctx, span := h.tracer.Start(ctx, "CreateAccountHandler.Handle")
	defer span.End()

	input := dto.CreateAccountInput{
		AccountExternalID: req.GetAccountExternalId(),
		AccountNumber:     req.GetAccountNumber(),
		TaxID:             req.GetTaxId(),
		Currency:          string(h.CurrencyTranslate(req)),
	}

	response, err := h.useCase.Execute(ctx, input)
	if err != nil {
		return nil, translator.ToGRPCError(err)
	}

	return &pb.CreateAccountResponse{
		AccountId: *response.AccountID,
	}, nil
}

func (h *CreateAccountHandler) CurrencyTranslate(req *pb.CreateAccountRequest) money.Currency {
	switch req.GetCurrency() {
	case pb.Currency_CURRENCY_BRL:
		return money.BRL
	case pb.Currency_CURRENCY_USD:
		return money.USD
	case pb.Currency_CURRENCY_EUR:
		return money.EUR
	default:
		return ""

	}
}
