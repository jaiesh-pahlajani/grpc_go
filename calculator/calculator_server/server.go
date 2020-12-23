package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc_go/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("function was inoked with %v", req)
	firstNumber := req.GetMessage().GetFirstNumber()
	secondNumber := req.GetMessage().GetLastNumber()

	sumResponse := &calculatorpb.SumResponse{
		Result: firstNumber + secondNumber,
	}

	return sumResponse, nil
}

func main() {

	// Port binding
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	// GRPC Server
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("Server started")

	// Bind port to grpc server
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
