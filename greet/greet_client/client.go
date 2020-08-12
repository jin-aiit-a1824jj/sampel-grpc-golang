package main

import (
	"context"
	"fmt"
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

	doUnary(c)
	doUnaryExercise(c)
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
