package decoder

import (
	"encoding/json"
	"net/http"
)

func ResponseSuccess[T any](write http.ResponseWriter, status int, data T) {
	write.Header().Set(ContentType, ApplicationJSON)
	write.WriteHeader(status)
	result := data

	if any(result) == nil {
		_ = json.NewEncoder(write)
		return
	}

	_ = json.NewEncoder(write).Encode(result)
}
