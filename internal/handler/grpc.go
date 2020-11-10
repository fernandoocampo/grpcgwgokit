package handler

import (
	"context"
	"errors"
	"log"

	"github.com/fernandoocampo/grpcgwgokit/internal/domain"
	"github.com/fernandoocampo/grpcgwgokit/internal/service"
	pb "github.com/fernandoocampo/grpcgwgokit/pkg/proto/grpcgwgokit/pb"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

// yourServiceServer implement reference data grpc server interface.
type yourServiceServer struct {
	echoHandler grpcTransport.Handler
	pb.UnimplementedYourServiceServer
}

// NewGRPCServer is a factory to create grpc servers for this project.
func NewGRPCServer(endpoints service.Endpoints) pb.YourServiceServer {
	return &yourServiceServer{
		echoHandler: grpcTransport.NewServer(
			endpoints.EchoEndpoint,
			DecodeGRPCEchoRequest,
			EncodeGRPCEchoResponse,
		),
	}
}

// CreateUser creates a new user. through grpc.
func (r *yourServiceServer) Echo(ctx context.Context, req *pb.StringMessage) (*pb.StringMessage, error) {
	_, resp, err := r.echoHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StringMessage), nil
}

// DecodeGRPCEchoRequest decodes grpc request for echo.
func DecodeGRPCEchoRequest(ctx context.Context, r interface{}) (interface{}, error) {
	log.Println("decoding grpc echo request")
	req, ok := r.(*pb.StringMessage)
	if !ok {
		log.Printf("decoding was not possible because request is not a pb.StringMessage: %+v", r)
		return nil, errors.New("request is not a pb.StringMessage")
	}
	newEcho := NewNewEcho(req).toEcho()
	log.Printf("new echo request was decoded successfully in %+v", newEcho)
	return newEcho, nil
}

// EncodeGRPCEchoResponse decode echo response.
func EncodeGRPCEchoResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res, ok := r.(domain.Echo)
	if !ok {
		log.Printf("decoding response was not possible because response is not a domain.Echo: %+v", r)
		return nil, errors.New("response is not a domain.Echo")
	}
	response := &pb.StringMessage{
		Value: res.Value,
	}

	return response, nil
}
