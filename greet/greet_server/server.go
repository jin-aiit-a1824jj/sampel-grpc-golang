package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

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
	fmt.Printf("result-> %v\n", result)
	res := &greetpb.CalculatorResponse{
		Result:      result,
		ResultInt64: firstNumber + secondNumber,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManytimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) PrimeNumber(req *greetpb.PrimeNumberRequest, stream greetpb.GreetService_PrimeNumberServer) error {
	fmt.Printf("PrimeNumber function was invoked with %v\n", req)
	n := req.GetNumber()
	k := int32(2)
	for i := 0; n > 1; i++ {
		if n%k == 0 {
			result := strconv.Itoa(i) + " -> [" + strconv.Itoa(int(k)) + "]"
			n = n / k
			res := &greetpb.PrimeNumberResponse{
				Result: result,
			}
			stream.Send(res)
			//time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request")

	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the cliend stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading cliend stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello" + firstName + "! "
	}
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
