package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/fernandoocampo/grpcgwgokit/pkg/proto/grpcgwgokit/pb"
)

func main() {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	// conn, err := grpc.Dial("localhost:8081", opts...)
	conn, err := grpc.Dial("localhost:50501", opts...)
	if err != nil {
		panic(err)
	}

	client := pb.NewYourServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	message := &pb.StringMessage{
		Value: "Hola",
	}
	answer, err := client.Echo(ctx, message)

	if err != nil {
		panic(err)
	}

	fmt.Println("Answer", answer)

	defer conn.Close()
}
