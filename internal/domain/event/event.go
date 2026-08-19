package event

import "context"

// Event define o contrato que todo evento de domínio precisa cumprir
type Event interface {
	SubjectName() string
	Payload() ([]byte, error)
}

// Publisher define o contrato para disparar eventos para o broker
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
