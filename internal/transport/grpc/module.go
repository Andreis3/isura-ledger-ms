package grpc

import "google.golang.org/grpc"

// Module defines the contract for registering handlers on the gRPC server.
type Module interface {
	Register(server *grpc.Server)
}
