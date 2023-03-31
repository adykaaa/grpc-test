package main

import (
	"context"
	"fmt"
	"log"

	"github.com/adykaaa/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial grpc server")
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	DoUnary(c)
}

func DoUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Toth",
			LastName:  "Adam",
		},
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		fmt.Printf("Could not send the request!")
	}
	log.Printf("repsonse from the server %v", resp.Result)

}
