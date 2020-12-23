package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc_go/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {

	// Create a connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)
	fmt.Printf("Client created: %v", client)

	message := &calculatorpb.Message{
		FirstNumber: 5,
		LastNumber:  10,
	}

	sumRequest := &calculatorpb.SumRequest{
		Message: message,
	}

	sumResponse, err := client.Sum(context.Background(), sumRequest)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sumResponse)
}
