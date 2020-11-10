package service

import (
	"context"
	"log"

	"github.com/fernandoocampo/grpcgwgokit/internal/domain"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints is a wrapper for endpoints
type Endpoints struct {
	EchoEndpoint endpoint.Endpoint
}

// NewEndpoints create an endpoints handler.
func NewEndpoints(service domain.EchoService) Endpoints {
	return Endpoints{
		EchoEndpoint: makeEchoEndpoint(service),
	}
}

// makeEchoEndpoint create endpoint for echo service.
func makeEchoEndpoint(service domain.EchoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		newecho := request.(*domain.Echo)

		result, err := service.DoEcho(ctx, newecho)
		if err != nil {
			log.Println("error", err)
		}
		return result, nil
	}
}
