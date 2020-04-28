package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nexlight101/gRPC_course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	// Create connection to the server
	options := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", options)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}

	// CLose the connection at exit
	defer cc.Close()

	// Establish a new calculator client
	cs := calculatorpb.NewCalculatorServiceClient(cc)
	fmt.Printf("Client activated: %v\n", cs)
	// send request to Greet unary client
	// doUnary(c)
	doSum(cs)
}

// doSum request the sum of two numbers
func doSum(cs calculatorpb.CalculatorServiceClient) {
	fmt.Println("Sending the Unary Sum request to server")
	req := &calculatorpb.CalculatorRequest{Sum: &calculatorpb.Sum{
		Number1: 3,
		Number2: 10,
	},
	}

	res, err := cs.Calculator(context.Background(), req)
	if err != nil {
		fmt.Printf("Error while calling Sum RPC: %v\n", err)
	}
	fmt.Printf("The sum of %d and %d is %d\n", req.Sum.Number1, req.Sum.Number2, res.GetResult())
}
