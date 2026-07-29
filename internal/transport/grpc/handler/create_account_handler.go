package handler

import (
	"context"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/application/command"
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

	input := command.CreateAccountInput{
		OwnerID:           req.GetOwnerId(),
		AccountExternalID: req.GetAccountNumber(),
		AccountNumber:     req.GetAccountNumber(),
		TaxID:             req.GetTaxId(),
		AccountingType:    h.AccountTypeTranslate(req),
		Currency:          h.CurrencyTranslate(req),
	}

	accountID, err := h.useCase.Execute(ctx, input)
	if err != nil {
		return nil, translator.ToGRPCError(err)
	}

	return &pb.CreateAccountResponse{
		AccountId: accountID,
	}, nil
}

func (h *CreateAccountHandler) AccountTypeTranslate(req *pb.CreateAccountRequest) string {
	switch req.GetAccountingType() {
	case pb.AccountingType_ACCOUNTING_TYPE_ASSET:
		return "ASSET"
	case pb.AccountingType_ACCOUNTING_TYPE_LIABILITY:
		return "LIABILITY"
	case pb.AccountingType_ACCOUNTING_TYPE_EQUITY:
		return "EQUITY"
	case pb.AccountingType_ACCOUNTING_TYPE_REVENUE:
		return "REVENUE"

	case pb.AccountingType_ACCOUNTING_TYPE_EXPENSE:
		return "EXPENSE"
	default:
		return ""

	}
}

func (h *CreateAccountHandler) CurrencyTranslate(req *pb.CreateAccountRequest) string {
	switch req.GetCurrency() {
	case pb.Currency_CURRENCY_BRL:
		return "BRL"
	case pb.Currency_CURRENCY_USD:
		return "USD"
	case pb.Currency_CURRENCY_EUR:
		return "EUR"
	default:
		return ""

	}
}
