package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"../greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("0.0.0.0:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Creted client:%f", c)

	//doUnary(c)
	//doUnaryExercise(c)
	//doServerStreaming(c)
	//doServerStreamingExercise(c)
	//doClientStreaming(c)
	//doClientStreamingExercise(c)
	//doBiDiStreaming(c)
	//doBiDiStreamingExercise(c)

	doUnaryWithDeadLine(c, 5*time.Second) //should complete
	doUnaryWithDeadLine(c, 1*time.Second) //should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "FirstName",
			LastName:  "LastName",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doUnaryExercise(c greetpb.GreetServiceClient) {
	req := &greetpb.CalculatorRequest{
		Calculator: &greetpb.Calculator{
			FirstNumber:  3.0,
			SecondNumber: 10.0,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculate RPC: %v", err)
	}
	log.Printf("Response from Calculate: %s", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "FirstName",
			LastName:  "LastName",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
		}
		log.Printf("Responese from GreetManyTimes: %v\n", msg.GetResult())
	}

}

func doServerStreamingExercise(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Server Streaming Exercise RPC...")

	req := &greetpb.PrimeNumberRequest{
		Number: int32(120),
	}
	resStream, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumber RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while calling PrimeNumber RPC: %v", err)
		}
		log.Printf("Responese from PrimeNumber: %v\n", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Client Streaming RPC...\n")

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 4",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 5",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range request {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving from LongGreet %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doClientStreamingExercise(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a Client Streaming RPC...\n")

	request := []*greetpb.ComputeAverageRequest{
		&greetpb.ComputeAverageRequest{
			Number: 1,
		},
		&greetpb.ComputeAverageRequest{
			Number: 2,
		},
		&greetpb.ComputeAverageRequest{
			Number: 3,
		},
		&greetpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range request {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		//time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving from ComputeAverage %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a BiDi Streaming RPC...\n")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while create stream: %v", err)
		return
	}

	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 1",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 4",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "FirstName - 5",
			},
		},
	}

	waitc := make(chan struct{})

	//we send a bunch of message to the client (go routine)
	go func() {
		// function to send a bunch of message
		for _, req := range request {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messsages from the client (go routine)
	go func() {
		// function to receive a bunch of message
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doBiDiStreamingExercise(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting to do a BiDi Streaming Exercise RPC...\n")

	// we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while create stream: %v", err)
		return
	}

	request := []*greetpb.FindMaximumRequest{
		&greetpb.FindMaximumRequest{
			Number: 1,
		},
		&greetpb.FindMaximumRequest{
			Number: 5,
		},
		&greetpb.FindMaximumRequest{
			Number: 3,
		},
		&greetpb.FindMaximumRequest{
			Number: 6,
		},
		&greetpb.FindMaximumRequest{
			Number: 2,
		},
		&greetpb.FindMaximumRequest{
			Number: 20,
		},
	}

	waitc := make(chan struct{})

	//we send a bunch of message to the client (go routine)
	go func() {
		// function to send a bunch of message
		for _, req := range request {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messsages from the client (go routine)
	go func() {
		// function to receive a bunch of message
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a doUnaryWithDeadLine RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "FirstName",
			LastName:  "LastName",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			}
		} else {
			log.Fatalf("error while calling GreetWithDeadline RPC: %v\n", statusErr)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
