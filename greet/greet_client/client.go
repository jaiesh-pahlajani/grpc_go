package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc_go/greet/greetpb"
	"google.golang.org/grpc"
)

func doUnary(client greetpb.GreetServiceClient) {
	greeting := &greetpb.Greeting{
		FirstName: "James",
		LastName:  "Bond",
	}

	greetRequest := &greetpb.GreetRequest{
		Greeting: greeting,
	}

	greetResponse, err := client.Greet(context.Background(), greetRequest)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(greetResponse)
}

func doServerStreaming(client greetpb.GreetServiceClient) {

	greeting := &greetpb.Greeting{
		FirstName: "Lionel",
		LastName:  "Messi",
	}

	greetRequest := &greetpb.GreetRequest{
		Greeting: greeting,
	}

	resStream, err := client.GreetManyTimes(context.Background(), greetRequest)
	if err != nil {
		fmt.Printf("error while calling server streaming %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we have reached end of stream
			break
		}
		if err != nil {
			fmt.Printf("Error %v", err)
		}
		fmt.Printf("Server streaming response %v \n", msg.GetResult())
	}
}

func doClientStreaming(client greetpb.GreetServiceClient) {

	greetingList := []*greetpb.Greeting{
		{
			FirstName: "James",
			LastName:  "Bond",
		},
		{
			FirstName: "Lionel",
			LastName:  "Messi",
		},
	}
	stream, err := client.LongGreet(context.Background())
	for _, greeting := range greetingList {
		req := &greetpb.LongGreetRequest{
			Greeting: greeting,
		}

		err = stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Sending request %v \n", req)
		if err != nil {
			fmt.Printf("Error %v", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("error %v", err)
	}
	fmt.Printf("Response is %v \n", response)
}

func doBiDiStreaming(client greetpb.GreetServiceClient) {

	// Create a stream
	stream, err := client.BidirectionalGreet(context.Background())
	if err != nil {
		fmt.Printf("Error %v \n", err)
	}

	waitc := make(chan struct{})

	greetingList := []*greetpb.Greeting{
		{
			FirstName: "James",
			LastName:  "Bond",
		},
		{
			FirstName: "Lionel",
			LastName:  "Messi",
		},
		{
			FirstName: "Miss",
			LastName:  "Moneypenny",
		},
		{
			FirstName: "Sachin",
			LastName:  "Tendulkar",
		},
		{
			FirstName: "Donald",
			LastName:  "Trump",
		},
	}

	// Send bunch of messages - go routine
	go func() {
		for _, greeting := range greetingList {
			req := &greetpb.GreetRequest{
				Greeting: greeting,
			}
			time.Sleep(1000 * time.Millisecond)
			fmt.Printf("Sending: %v \n", greeting)
			err := stream.Send(req)
			if err != nil {
				fmt.Printf("Error %v \n", err)
			}
		}
		err := stream.CloseSend()
		if err != nil {
			fmt.Printf("Error %v \n", err)
		}
	}()

	// Receive bunch of messages do routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				fmt.Printf("Error %v \n", err)
			}
			fmt.Printf("Response stream: %v \n", res.GetResult())
		}
	}()

	<-waitc

	// Block until everything is done

}

func main() {

	// Create a connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Client created: %v \n", client)

	//doUnary(client)
	//doServerStreaming(client)
	//doClientStreaming(client)
	doBiDiStreaming(client)
}
