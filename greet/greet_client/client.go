package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
	doServerStreaming(c)
	doServerStreamingExercise(c)
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
