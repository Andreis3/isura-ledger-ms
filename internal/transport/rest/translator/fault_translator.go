package translator

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

type ProtocolError struct {
	HTTPStatus int
}

var TranslatorStatusCode = map[fault.Code]ProtocolError{
	fault.CodeBadRequest:          {HTTPStatus: http.StatusBadRequest},
	fault.CodeUnauthorized:        {HTTPStatus: http.StatusUnauthorized},
	fault.CodeForbidden:           {HTTPStatus: http.StatusForbidden},
	fault.CodeNotFound:            {HTTPStatus: http.StatusNotFound},
	fault.CodeConflict:            {HTTPStatus: http.StatusConflict},
	fault.CodeUnprocessableEntity: {HTTPStatus: http.StatusUnprocessableEntity},
	fault.CodeInternal:            {HTTPStatus: http.StatusInternalServerError},
}
