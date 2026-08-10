package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
	"github.com/andreis3/isura-ledger-ms/internal/transport/rest/translator"
)

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"
)

type TypeResponseError struct {
	CodeError       string         `json:"code_error"`
	Cause           string         `json:"cause,omitempty"`
	ErrorFields     map[string]any `json:"error_fields,omitempty"`
	FriendlyMessage any            `json:"friendly_message"`
}

func ResponseError(write http.ResponseWriter, err error) {
	if t, ok := errors.AsType[*fault.DomainError](err); ok {
		result := TypeResponseError{
			CodeError:       string(t.Code),
			Cause:           t.Cause.Error(),
			ErrorFields:     t.Fields,
			FriendlyMessage: t.FriendlyMessage,
		}

		write.Header().Set(ContentType, ApplicationJSON)
		write.WriteHeader(translator.TranslatorStatusCode[t.Code].HTTPStatus)
		_ = json.NewEncoder(write).Encode(result)
		return
	}

	write.Header().Set(ContentType, ApplicationJSON)
	write.WriteHeader(http.StatusInternalServerError)

	result := TypeResponseError{
		FriendlyMessage: "Internal server error",
	}

	_ = json.NewEncoder(write).Encode(result)
}
