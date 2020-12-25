package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc_go/greet/greetpb"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greetResponse := &greetpb.GreetResponse{
		Result: "Hey " + firstName + " " + lastName + "!",
	}
	return greetResponse, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Many Times function was invoked %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetResponse{
			Result: "hello " + firstName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil

}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Greet Many Times function was invoked")
	responseMessage := "Hey "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &greetpb.LongGreetResponse{
				Result: responseMessage,
			}
			err := stream.SendAndClose(response)
			if err != nil {
				fmt.Printf("Error %v", err)
			}
			return err
		}
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		responseMessage += firstName + " " + lastName + "!"
	}
}

func (*server) BidirectionalGreet(stream greetpb.GreetService_BidirectionalGreetServer) error {
	fmt.Println("Invoked Bidirectional Greet!")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client streaming %v \n", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "!"
		response := &greetpb.GreetResponse{
			Result: result,
		}
		err = stream.Send(response)
		if err != nil {
			fmt.Printf("Error %v \n", err)
			return err
		}
	}
	return nil
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet Deadline function was invoked %v \n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// Client cancelled the request
			fmt.Println("Client cancelled the request")
			return nil, status.Error(codes.Canceled, "the client cancelled the requesr")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greetResponse := &greetpb.GreetResponse{
		Result: "Hey " + firstName + " " + lastName + "!",
	}
	return greetResponse, nil
}

func main() {

	// Port binding
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	// SSL Cert
	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	// GRPC Server
	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Register reflection server on grpc
	reflection.Register(s)

	fmt.Println("Server started")

	// Bind port to grpc server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
