package service

import (
	"context"

	"github.com/fernandoocampo/grpcgwgokit/internal/domain"
)

const silence = "zzzzzzzzzz"

type basicEcho struct {
}

// NewBasicEcho creates a basic implementation for echo service.
func NewBasicEcho() domain.EchoService {
	return &basicEcho{}
}

// DoEcho does a big echo
func (b *basicEcho) DoEcho(ctx context.Context, echo *domain.Echo) (domain.Echo, error) {
	value := silence
	var result domain.Echo
	if echo != nil && echo.Value != "" {
		value = echo.Value
	}

	result.Value = value

	return result, nil
}
