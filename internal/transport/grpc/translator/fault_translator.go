package translator

import (
	"errors"
	"fmt"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProtocolError mapeia um Code de domínio para status de protocolos externos.
// Adicione novos protocolos aqui (gRPC, GraphQL, etc.) sem tocar no domínio.
type ProtocolError struct {
	GRPCCode codes.Code
}

// translator é o mapa interno de conversão Code → protocolos.
var translator = map[fault.Code]ProtocolError{
	fault.CodeBadRequest:           {GRPCCode: codes.InvalidArgument},
	fault.CodeUnauthorized:         {GRPCCode: codes.Unauthenticated},
	fault.CodeForbidden:            {GRPCCode: codes.PermissionDenied},
	fault.CodeNotFound:             {GRPCCode: codes.NotFound},
	fault.CodeConflict:             {GRPCCode: codes.AlreadyExists},
	fault.CodeUnprocessableEntity:  {GRPCCode: codes.InvalidArgument},
	fault.CodeInternal:             {GRPCCode: codes.Internal},
	fault.CodeDatabaseError:        {GRPCCode: codes.Internal},
	fault.CodeInvalidEntity:        {GRPCCode: codes.InvalidArgument},
	fault.CodeUnknown:              {GRPCCode: codes.Unknown},
	fault.CodeCacheError:           {GRPCCode: codes.Internal},
	fault.CodeExternalService:      {GRPCCode: codes.Unavailable},
	fault.CodeTimeoutError:         {GRPCCode: codes.DeadlineExceeded},
	fault.CodeInvalidTransfer:      {GRPCCode: codes.InvalidArgument},
	fault.CodeInsufficientBalance:  {GRPCCode: codes.InvalidArgument},
	fault.CodeDuplicateTransaction: {GRPCCode: codes.InvalidArgument},
}

// GRPCStatus retorna o status GRPC correspondente ao erro.
// Se o erro não for um DomainError, retorna 500.
// Se o Code não estiver mapeado, retorna 500.
func GRPCStatus(err error) codes.Code {
	var de *fault.DomainError
	if !errors.As(err, &de) {
		return codes.Internal
	}

	if p, ok := translator[de.Code]; ok {
		return p.GRPCCode
	}

	return codes.Internal
}

// Response é a estrutura que vai no body da resposta de erro HTTP.
// Nunca exponha DomainError.Error() aqui — contém informação técnica.
type Response struct {
	Code    fault.Code     `json:"code"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

// ToGRPCError converte um DomainError para um erro gRPC com detalhes customizados.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Converte para DomainError
	if _, ok := errors.AsType[*fault.DomainError](err); !ok {
		// Se não for DomainError, retorna erro interno genérico
		return status.Error(codes.Internal, "internal server error")
	}

	// Mapeia o código do domínio para código gRPC
	grpcCode := GRPCStatus(err) // sua função já existe

	// Cria a resposta amigável
	resp := ToResponse(err) // retorna Response com Code, Message, Fields

	var br *errdetails.BadRequest
	// Cria o status base com a mensagem amigável
	st := status.New(grpcCode, resp.Message)
	// Preenche os FieldError se houver campos
	if len(resp.Fields) > 0 {
		fields := make([]*errdetails.BadRequest_FieldViolation, 0, len(resp.Fields))
		for field, val := range resp.Fields {
			// Converte o valor para string (se for erro de validação)
			msg := fmt.Sprintf("%v", val)
			fields = append(fields, &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: msg,
			})
		}
		br = &errdetails.BadRequest{FieldViolations: fields}
		// Adiciona o ErrorResponse como detalhe
		stWithDetails, err := st.WithDetails(br)
		if err != nil {
			// Se falhar, retorna o status sem detalhes (mas ainda com a mensagem)
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return st.Err()
}

// ToResponse converte um DomainError para a resposta segura ao client.
func ToResponse(err error) Response {
	var de *fault.DomainError
	if !errors.As(err, &de) {
		return Response{
			Code:    fault.CodeInternal,
			Message: "Internal server error",
		}
	}

	return Response{
		Code:    de.Code,
		Message: de.FriendlyMessage,
		Fields:  de.Fields,
	}
}
