package model

import "fmt"

// ErrorKind lets the transport layer map domain failures to stable HTTP status codes.
type ErrorKind string

const (
	ErrorInvalid  ErrorKind = "invalid"
	ErrorNotFound ErrorKind = "not_found"
	ErrorConflict ErrorKind = "conflict"
	ErrorInternal ErrorKind = "internal"
)

// DomainError preserves a machine-readable category without exposing storage details.
type DomainError struct {
	Kind    ErrorKind
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewError(kind ErrorKind, format string, args ...any) error {
	return &DomainError{Kind: kind, Message: fmt.Sprintf(format, args...)}
}
