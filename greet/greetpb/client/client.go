package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/adykaaa/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func DoUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Toth",
			LastName:  "Adam",
		},
	}

	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		err, ok := status.FromError(err)
		if ok  {
			if err.Code() == codes.InvalidArgument {
				log.Fatal("Cannot providte an empty FirstName")
			}
	}
	log.Printf("repsonse from the server %v", resp.Result)
}

func DoServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting server streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Toth",
			LastName:  "Adam",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			//we've reached end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}

		log.Printf("Response: " + msg.GetResult())
	}
}

func DoClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Adam",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "asadsadsa",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "LOLOLOLOL",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Madeupname",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling Long Greet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request %v \n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error during close and recieve! %v", err)
	}

	fmt.Printf("server response after client stream: %v", resp)
}

func DoBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming RPC...")
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Adam",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "asadsadsa",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "LOLOLOLOL",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Madeupname",
			},
		},
	}

	//we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitch := make(chan struct{})
	//we send bunch of messages to the server
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message %v \n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()
	//we receive a bunch of messages from the server
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("Error while receving %v", err)
				break
			}
			fmt.Printf("Received: %v \n", resp.GetResult())
		}
		close(waitch)
	}()
	//block until everything is done
	<-waitch
}
func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial grpc server")
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//DoUnary(c)
	//DoClientStreaming(c)
	//DoServerStreaming(c)
	DoBiDiStreaming(c)
}
