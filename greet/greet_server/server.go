package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

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
	fmt.Println("LongGreet function was invoked with a streaming request")

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

func (*server) ComputeAverage(stream greetpb.GreetService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage function was invoked with a streaming request")

	result := int64(0)
	for i := 0; ; i++ {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the cliend stream
			return stream.SendAndClose(&greetpb.ComputeAverageResponse{
				Result: "result -> " + strconv.FormatFloat(float64(int(result))/float64(i), 'f', 2, 64),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading cliend stream: %v", err)
		}
		result += req.GetNumber()
	}

}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function was invoked with a streaming request")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Errorr while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*server) FindMaximum(stream greetpb.GreetService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request")

	numberSlice := []int{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Errorr while reading client stream: %v", err)
			return err
		}

		numberSlice = append(numberSlice, int(req.GetNumber()))
		sort.Slice(numberSlice, func(i, j int) bool {
			return numberSlice[i] < numberSlice[j]
		})

		ans := numberSlice[len(numberSlice)-1]
		result := "FindMaximum [" + strconv.Itoa(ans) + "] "
		fmt.Println("Send response ->" + result)
		sendErr := stream.Send(&greetpb.FindMaximumResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
