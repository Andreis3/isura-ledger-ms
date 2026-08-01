package middleware

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/application"
)

// responseWriterWrapper intercepta o status HTTP para podermos registrar no span se necessário
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterWrapper) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Tracing(tracer application.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Usa o camimnho  da URL (ex: /accounts) como nome do span
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(r.Context(), spanName)
			defer span.End()

			// Injeta o contexto atualizado com span na requisição
			r = r.WithContext(ctx)

			// Wrapper para capturar o status code da resposta
			ww := &responseWriterWrapper{w, http.StatusOK}

			// Executa o próximo handler
			next.ServeHTTP(ww, r)

			// Se o status for de erro (>= 500), podemos registrar como erro no span
			if ww.statusCode >= http.StatusInternalServerError {
				span.RecordError(http.ErrAbortHandler) // Ou pode registrar uma mensagem/erro customizado
			}
		})
	}
}
