package domain

import "context"

// Echo contains data of a echo.
type Echo struct {
	Value string
}

// EchoService defines behavior for echo service
type EchoService interface {
	// DoEcho does a big echo
	DoEcho(ctx context.Context, echo *Echo) (Echo, error)
}
