package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/adykaaa/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet function was invoked with %v", req)
	firstName := req.Greeting.FirstName
	return &greetpb.GreetResponse{Result: "Hello " + firstName}, nil
}

func (*server) mustEmbedUnimplementedGreetServiceServer() {}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("cannot start listening")
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
