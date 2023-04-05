package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/adykaaa/grpc-test/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

// unary
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet function was invoked with %v \n", req)
	if len(req.Greeting.FirstName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "The FirstName cannot be empty!")
	}
	return &greetpb.GreetResponse{Result: "Hello " + req.Greeting.FirstName + req.Greeting.LastName}, nil
}

// server streaming
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)

	for i := 0; i < 10; i++ {
		result := "Hello " + req.Greeting.FirstName + req.Greeting.LastName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// client streaming
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("GreetManyTimes function was invoked with a client streaming request \n")

	result := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.Greeting.FirstName
		result += "Hello " + firstName + "! "
		log.Print(result)
	}
}

// bidirectional streaming
func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked \n")
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Fatalf("Errorr while reading client stream: %v", err)
			return err
		}
		firstName := req.Greeting.FirstName
		result := "Hello " + firstName + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalf("Errorr while sending data to client: %v", err)
			return err
		}

	}
}

func (*server) mustEmbedUnimplementedGreetServiceServer() {}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("cannot start listening")
	}
	log.Print("Server is listening! \n")

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
