package translator

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

type ProtocolError struct {
	HTTPStatus int
}

var TranslatorStatusCode = map[fault.Code]ProtocolError{
	fault.CodeBadRequest:           {HTTPStatus: http.StatusBadRequest},
	fault.CodeUnauthorized:         {HTTPStatus: http.StatusUnauthorized},
	fault.CodeForbidden:            {HTTPStatus: http.StatusForbidden},
	fault.CodeNotFound:             {HTTPStatus: http.StatusNotFound},
	fault.CodeConflict:             {HTTPStatus: http.StatusConflict},
	fault.CodeUnprocessableEntity:  {HTTPStatus: http.StatusUnprocessableEntity},
	fault.CodeInternal:             {HTTPStatus: http.StatusInternalServerError},
	fault.CodeUnknown:              {HTTPStatus: http.StatusInternalServerError},
	fault.CodeDatabaseError:        {HTTPStatus: http.StatusInternalServerError},
	fault.CodeCacheError:           {HTTPStatus: http.StatusInternalServerError},
	fault.CodeExternalService:      {HTTPStatus: http.StatusBadGateway},
	fault.CodeTimeoutError:         {HTTPStatus: http.StatusGatewayTimeout},
	fault.CodeInvalidEntity:        {HTTPStatus: http.StatusBadRequest},
	fault.CodeInvalidTransfer:      {HTTPStatus: http.StatusBadRequest},
	fault.CodeInsufficientBalance:  {HTTPStatus: http.StatusBadRequest},
	fault.CodeDuplicateTransaction: {HTTPStatus: http.StatusBadRequest},
}
