package middleware

import (
	"net/http"
	"time"

	"github.com/andreis3/isura-ledger-ms/internal/application"
)

type responseWriterWrapperMetrics struct {
	http.ResponseWriter
	statusCode int
}

// MetricsMiddleware é um middleware HTTP que registra métricas de requisições, equivalente ao MetricsInterceptor do gRPC.
func MetricsMiddleware(metrics application.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrapper para capturar o status HTTP gerado pelo handler
			ww := &responseWriterWrapperMetrics{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // padrão do net/http se WriteHeader não for chamado explicitamente
			}

			// Executa o próximo handler da cadeia
			next.ServeHTTP(ww, r)

			duration := float64(time.Since(start).Milliseconds())
			statusCode := ww.statusCode

			// Usamos o padrão r.Method + " " + r.URL.Path (ou r.Pattern se estiver usando Go 1.22+ routing nativo)
			// como identificador da rota, de forma análoga ao info.FullMethod do gRPC.
			routeIdentifier := r.Method + " " + r.URL.Path

			metrics.RecordRequestTotal(routeIdentifier, "http", statusCode)
			metrics.RecordRequestDuration(routeIdentifier, "http", statusCode, duration)
		})

	}
}
