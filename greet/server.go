package main

import (
	"context"

	"github.com/adykaaa/grpc-test/greet/greetpb"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

}
