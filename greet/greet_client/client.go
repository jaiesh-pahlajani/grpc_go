package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doServerStreaming(client)
}
