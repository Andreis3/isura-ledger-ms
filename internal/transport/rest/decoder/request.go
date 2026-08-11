package decoder

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

func RequestDecoder[T any](req *http.Request) (T, error) {
	defer req.Body.Close()
	var result T

	// Use the Stream Decoder directly on the io.Reader of the request body.
	// This avoids allocating a giant slice on the heap with io.ReadAll and uses Sonic's internal pool.
	err := jsonEngine.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		return result, fault.ErrorJSON(err)
	}

	return result, nil
}