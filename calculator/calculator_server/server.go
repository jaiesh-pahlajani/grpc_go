package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

func (*server) PrimeNumberDecomposition(numberRequest *calculatorpb.NumberRequest, decompositionServer calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	fmt.Printf("Prime number decomposition was inoked with %v \n", numberRequest)

	n := numberRequest.GetNumber()

	k := int32(2)
	for n > 1 {
		if n%k == 0 {
			n = n / k
			numberResponse := &calculatorpb.NumberResponse{
				Number: k,
			}
			time.Sleep(1000 * time.Millisecond)
			err := decompositionServer.Send(numberResponse)
			if err != nil {
				fmt.Printf("Error while sending to client %v", err)
			}
		} else {
			k++
		}
	}

	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Println("Received request")
	numberResponse := &calculatorpb.AverageNumberResponse{
		Number: 0,
	}
	ctr := 0
	for {
		numberRequest, err := stream.Recv()
		if err == io.EOF {
			numberResponse.Number /= float32(ctr)
			err := stream.SendAndClose(numberResponse)
			if err != nil {
				fmt.Printf("Error %v \n", err)
			}
			return err

		}
		if err != nil {
			fmt.Printf("Error %v \n", err)
		}
		ctr++
		numberResponse.Number += float32(numberRequest.Number)
		fmt.Printf("sum %v \n", numberResponse.Number)
	}
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
