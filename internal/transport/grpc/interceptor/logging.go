package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor is a unary interceptor that logs each gRPC request.
// It must be registered in grpc.NewServer via grpc.UnaryInterceptor.
func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		// calls the real handler

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.String("status", code.String()),
			slog.String("duration", duration.String()),
		}

		// codes that are client errors — logs as WARN
		// codes that are server errors — logs as ERROR
		switch code {
		case codes.Internal, codes.Unavailable, codes.DataLoss, codes.Unknown:
			log.ErrorContext(ctx, "gRPC request failed", attrs...)
		case codes.OK:
			log.InfoContext(ctx, "gRPC request completed", attrs...)
		default:
			log.WarnContext(ctx, "gRPC request rejected", attrs...)
		}

		return resp, err
	}
}
