package decoder

import (
	"net/http"

	"github.com/bytedance/sonic"
)

func ResponseSuccess[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonBytes, err := sonic.Marshal(data)
	if err != nil {
		// Tratar erro de serialização se necessário
		return
	}
	_, _ = w.Write(jsonBytes)
}
