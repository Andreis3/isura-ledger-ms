package middleware

import (
	"net/http"

	"github.com/andreis3/isura-ledger-ms/internal/application"
)

// responseWriterWrapper intercepts the HTTP status so we can record it in the span if necessary.
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
			// Uses the URL path (e.g., /accounts) as the span name.
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(r.Context(), spanName)
			defer span.End()

			// Injects the updated context with the span into the request.
			r = r.WithContext(ctx)

			// Wrapper to capture the response status code.
			ww := &responseWriterWrapper{w, http.StatusOK}

			// Executes the next handler.
			next.ServeHTTP(ww, r)

			// If the status is an error (>= 500), we can record it as an error in the span.
			if ww.statusCode >= http.StatusInternalServerError {
				span.RecordError(http.ErrAbortHandler) // Or you can record a custom message/error.
			}
		})
	}
}
