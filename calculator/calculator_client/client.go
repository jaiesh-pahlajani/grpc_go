package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc_go/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func doUnary(client calculatorpb.CalculatorServiceClient) {
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

func doServerStreaming(client calculatorpb.CalculatorServiceClient) {
	numberRequest := &calculatorpb.NumberRequest{
		Number: 120,
	}

	streamRes, err := client.PrimeNumberDecomposition(context.Background(), numberRequest)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	for {
		numberResponse, err := streamRes.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		fmt.Printf("Server streaming response %v \n", numberResponse.GetNumber())
	}

}

func doClientStreaming(client calculatorpb.CalculatorServiceClient) {
	n := []int{1, 3, 4, 5, 5}
	stream, err := client.Average(context.Background())
	if err != nil {
		fmt.Printf("Error %v", err)
	}
	for _, i := range n {
		req := &calculatorpb.NumberRequest{
			Number: int32(i),
		}
		err := stream.Send(req)
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		fmt.Printf("Sent %v \n", i)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("Error %v", err)
	}
	fmt.Printf("Response is %v \n", res)
}

func doBiDiSteaming(client calculatorpb.CalculatorServiceClient) {

	n := []int{1, 3, 10, 4, 5, 5}
	stream, err := client.CurrentMax(context.Background())
	if err != nil {
		fmt.Printf("Error %v", err)
	}
	waitc := make(chan struct{})

	// Send
	go func() {
		for _, i := range n {
			req := &calculatorpb.NumberRequest{
				Number: int32(i),
			}
			time.Sleep(1000 * time.Millisecond)
			fmt.Printf("Stream send %v \n", req.Number)
			err := stream.Send(req)
			if err != nil {
				fmt.Printf("Error %v", err)
			}
		}
		err := stream.CloseSend()
		if err != nil {
			fmt.Printf("Error %v", err)
		}
	}()

	// Recv
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("Error %v", err)
			}
			res.GetNumber()
			fmt.Printf("Stream Response %v \n", res.GetNumber())
		}
		close(waitc)
	}()

	// block
	<-waitc
}

func main() {

	// Create a connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)
	fmt.Printf("Client created: %v", client)

	//doUnary(client)
	//doServerStreaming(client)
	//doClientStreaming(client)
	doBiDiSteaming(client)
}
