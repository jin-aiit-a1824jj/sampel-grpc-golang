package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"../greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) Calculate(ctx context.Context, req *greetpb.CalculatorRequest) (*greetpb.CalculatorResponse, error) {
	fmt.Printf("Calculate function was invoked with %v\n", req)
	firstNumber := req.GetCalculator().GetFirstNumber()
	secondNumber := req.GetCalculator().GetSecondNumber()
	result := "Calculate -> " + strconv.FormatInt(firstNumber, 10) + " + " + strconv.FormatInt(secondNumber, 10) + " = " + strconv.FormatInt(firstNumber+secondNumber, 10)
	fmt.Printf("result-> %v", result)
	res := &greetpb.CalculatorResponse{
		Result:      result,
		ResultInt64: firstNumber + secondNumber,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
