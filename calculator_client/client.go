package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doSum(cs)

	// send a request to Prime number API
	// doPrime(cs)

	//Send request to GCF API
	// doGCF(cs)

	// Send a sentance to Split
	// doBWord(cs)

	// send a word to split
	// doLetters(cs)

	// doAverage sends a stream of numbers to the server to average
	// doAverage(cs)

	// doMultiply sends a stream of numbers to the server to multiply
	// doMultiply(cs)

	// doFindMaximum sends a stream of numbers and receives a stream of maximums
	doFindMaximum(cs)

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

func doPrime(cs calculatorpb.CalculatorServiceClient) {
	fmt.Println("Sending the Prime Factor request to server")
	req := &calculatorpb.PrimeNumberRequest{
		Number: 120,
	}
	res, err := cs.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while receiving from Prime RPC: %v\n", err)
	}
	for {
		response, rErr := res.Recv()
		if rErr == io.EOF {
			break
		}
		fmt.Printf("A Prime Factor of %v is: %v\n", req.GetNumber(), response.GetResult())
	}
}

func doGCF(cs calculatorpb.CalculatorServiceClient) {
	fmt.Println("Sending the GCF request to server")
	req := &calculatorpb.GCFRequest{
		Number1: 24,
		Number2: 60,
	}
	res, err := cs.GCF(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while receiving from GCF RPC: %v\n", err)
	}
	fmt.Printf("The Greatest Common Factor of %v and %v: ", req.GetNumber1(), req.GetNumber2())
	product := 1
	for {
		response, rErr := res.Recv()
		if rErr == io.EOF {
			break
		}
		fmt.Printf(" %v ", response.GetResult())
		product *= int(response.GetResult())
	}
	fmt.Println("=", product)
}

// doBWord Splits a sentance into single words
func doBWord(cs calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.BreakWordRequest{
		Word: "Hello there you turd",
	}
	fmt.Printf("Sending the sentance '%v' to split request to server\n", req.GetWord())
	res, err := cs.BreakWord(context.Background(), req)

	if err != nil {
		log.Fatalf("Cannot communicate with RPC %v\n", err)
	}
	for {
		response, rErr := res.Recv()
		if rErr == io.EOF {
			break
		}
		fmt.Printf("  '%v'  ", response.GetResult())
	}

}

// function doLetters() calls a letters API
func doLetters(cs calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.LettersRequest{
		Word: "Howdy there partner?",
	}
	fmt.Printf("Sending the sentance '%v' to split into letters request to server\n", req.GetWord())
	res, err := cs.Letters(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not reveive response: %v\n", err)
	}
	for {
		response, rErr := res.Recv()
		if rErr == io.EOF {
			break
		}
		fmt.Printf(" %v ", response.GetResult())
	}
}

func doAverage(cs calculatorpb.CalculatorServiceClient) {
	fmt.Println("Sending numbers to server")
	cl, err := cs.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Could not get response: %v\n", err)
	}
	for i := 1; i < 5; i++ {
		req := &calculatorpb.ComputeAverageRequest{
			Number: int32(i),
		}
		err1 := cl.Send(req)
		fmt.Println(req.GetNumber())
		if err1 != nil {
			log.Fatalf("Could not get response: %v\n", err)
		}
	}
	res, rErr := cl.CloseAndRecv()
	if rErr != nil {
		log.Fatalf("Could not get response: %v\n", err)
	}
	fmt.Printf("The avarage of 1, 2, 3, 4: %.2f", res.GetResult())
}

// doMultiply streams numbers to multiply
func doMultiply(cs calculatorpb.CalculatorServiceClient) {

	reqX := []*calculatorpb.MultiplierRequest{
		{
			Number: 6,
		},
		{
			Number: 5,
		},
		{
			Number: 4,
		},
		{
			Number: 7,
		},
	}
	client, err := cs.Multiplier(context.Background())
	if err != nil {
		log.Fatalf("Could not communicate with server : %v\n", err)
	}
	fmt.Println("Mutlipling:")
	for _, v := range reqX {
		sErr := client.Send(v)
		if sErr != nil {
			log.Fatalf("Could not send streaming request to the server : %v\n", sErr)
		}
		fmt.Printf(" %v ", v)
	}
	response, rErr := client.CloseAndRecv()
	if rErr != nil {
		log.Fatalf("Did not receive response from server: %v\n", rErr)
	}
	fmt.Printf("\nThe product is: %v", response.GetResult())
}

// doFindMaximum receives stream numbers and sends back a stream of maximums

func doFindMaximum(cs calculatorpb.CalculatorServiceClient) {
	fmt.Println("Sending numbers to server")
	waitc := make(chan struct{})
	// Create a stream from client
	stream, sErr := cs.FindMaximum(context.Background())
	if sErr == io.EOF {
		close(waitc)
	}
	if sErr != nil {
		log.Fatalf("Could not connect to stream: %v\n", sErr)
	}

	// Send a bunch of numbers
	go func() {
		numbersX := []int32{1, 5, 3, 6, 2, 20}
		for _, v := range numbersX {
			fmt.Printf("Sending number: %d ", v)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: v,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		cErr := stream.CloseSend()
		if cErr != nil {
			log.Fatalf("Could not close stream: %v\n", cErr)
		}
	}()

	// Receive a bunch of maximums
	go func() {
		for {
			res, rErr := stream.Recv()
			if rErr == io.EOF {
				break
			}
			if rErr != nil {
				log.Fatalf("Could not receive stream: %v\n", rErr)
				break
			}
			fmt.Printf(" Maximum = %d\n", res.GetResult())
		}
		close(waitc)
	}()

	// Block channel till done
	<-waitc
}
