package decoder

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

func RequestDecoder[T any](req *http.Request) (T, error) {
	defer req.Body.Close()
	var result T
	var jsonUnmarshalTypeError *json.UnmarshalTypeError
	var jsonSyntaxError *json.SyntaxError
	err := json.NewDecoder(req.Body).Decode(&result)
	switch {
	case errors.As(err, &jsonSyntaxError):
		return result, fault.ErrorJSONSyntaxError(jsonSyntaxError)

	case errors.As(err, &jsonUnmarshalTypeError):
		return result, fault.ErrorJSONUnmarshalTypeError(jsonUnmarshalTypeError)

	case err != nil:
		return result, fault.ErrorJSON(err)
	}

	return result, nil
}
