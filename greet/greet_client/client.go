package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpc_go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	// Create a connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer conn.Close()

	client := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Client created: %v", client)

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
