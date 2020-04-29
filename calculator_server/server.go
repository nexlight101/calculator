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
	fmt.Println("Hello from gRPC-Prime-Number-Server")

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

// PrimeNumber method retuns stream of prime factors
func (*server) PrimeNumber(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberServer) error {
	fmt.Printf("PrimeNumber Request received in server %v", req)
	resultX := toPrime(req.GetNumber())
	fmt.Println(resultX)
	for _, v := range resultX {
		stream.Send(&calculatorpb.PrimeNumberResponse{
			Result: v,
		})
	}
	return nil
}

// toPrime breaks down a number into prime factors
func toPrime(num int32) []int32 {
	// implement algrorithm
	primeX := make([]int32, 0, 20) // record all prime factors
	pFactor := 2
	n := int(num)
	for {
		if n <= 1 {
			return primeX
		}
		if n%pFactor == 0 { // if pFactor evenly divides into n
			primeX = append(primeX, int32(pFactor)) //add to the slice
			n = n / pFactor                         //remove prime factor from n
		} else {
			pFactor++
		}
	}
}

// Greatest Common Factor generator
func (*server) GCF(req *calculatorpb.GCFRequest, stream calculatorpb.CalculatorService_GCFServer) error {
	fmt.Printf("GCF Request received in server %v \n", req)
	// Get the numbers from the request
	number1 := req.GetNumber1()
	number2 := req.GetNumber2()
	// Make empty slices
	number1X := make([]int32, 0, 20)
	number2X := make([]int32, 0, 20)
	// Populate both slices
	number1X = toPrime(number1)
	number2X = toPrime(number2)
	// Determine common factors
	gcfX := make([]int32, 0, 20)
	d := 0
	for x, v := range number1X {
		fmt.Printf("outer loop index %v value %v \n", x, v)
		for i := d; i < len(number2X); {
			if number2X[i] == v {
				d = i + 1
				gcfX = append(gcfX, v)
				fmt.Printf("index %v ", i)
				fmt.Printf("value %v \n", v)
				break
			}
			i++
		}
	}
	for _, g := range gcfX {
		stream.Send(&calculatorpb.GCFResponse{
			Result: g,
		},
		)
		fmt.Println(g)
	}

	return nil
}
