package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"../greetpb"
	"google.golang.org/grpc"
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
	doClientStreaming(c)
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
