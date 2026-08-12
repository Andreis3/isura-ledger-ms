package entity

import (
	"errors"

	"github.com/google/uuid"
)

// ID represents a domain entity identifier.
type ID struct {
	value string
}

// NewID creates a new validated ID.
func NewID(value string) (ID, error) {
	if _, err := uuid.Parse(value); err != nil {
		return ID{}, errors.New("invalid ID format")
	}
	return ID{value: value}, nil
}

// NewIDV7 generates a new ID using UUIDv7.
func NewIDV7() (ID, error) {
	val, err := uuid.NewV7()
	if err != nil {
		return ID{}, err
	}
	return ID{value: val.String()}, nil
}

// String returns the string representation of the ID.
func (i ID) String() string {
	return i.value
}
