package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/nexlight101/gRPC_course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

// Create a server type
type server struct{}

// Sum method returns the sum of the requested numbers
func (*server) Calculator(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Sum Request received in server %v", req)
	result := req.GetSum().GetNumber1() + req.GetSum().GetNumber2()
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello from gRPC-Server")

	// Create listener
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//Create a new gRPC server
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Check if the server is serving the listener
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
