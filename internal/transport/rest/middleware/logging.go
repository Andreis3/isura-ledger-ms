package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriterWrapper intercepta o status HTTP para conseguirmos logar o código de resposta correto
type responseWriterWrapperLogging struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapperLogging) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Logging é um middleware HTTP que loga cada request, equivalente ao LoggingInterceptor do gRPC.
func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrapper para capturar o status HTTP gerado pelo handler
			ww := &responseWriterWrapperLogging{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // padrão do net/http se WriteHeader não for chamado explicitamente
			}

			// Executa o próximo handler da cadeia
			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			attrs := []any{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status_code", ww.statusCode),
				slog.String("duration", duration.String()),
			}

			// Mapeamento de severidade baseado no status code HTTP (análogo aos gRPC codes)
			switch {
			case ww.statusCode >= http.StatusInternalServerError: // 5xx: Erros de Servidor (Internal, Unavailable, etc.)
				log.ErrorContext(r.Context(), "HTTP request failed", attrs...)
			case ww.statusCode >= http.StatusBadRequest: // 4xx: Erros de Cliente (InvalidArgument, NotFound, etc.)
				log.WarnContext(r.Context(), "HTTP request rejected", attrs...)
			default: // 2xx / 3xx: Sucesso ou Redirecionamento
				log.InfoContext(r.Context(), "HTTP request completed", attrs...)
			}
		})
	}
}
