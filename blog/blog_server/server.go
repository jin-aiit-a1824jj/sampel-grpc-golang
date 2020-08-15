package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	"../blogpb"
)

var collection *mongo.Collection

type server struct {

}

type blogItem struct {
	ID 		 primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string 		    `bson:"author_id"`
	Content  string			    `bson:"content"`
	Title 	 string	            `bson:"title"`
}

func main() {
	// if we crash the go code , we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")

	// connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.20.0.20:27012"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Blog Service Started")
	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func(){
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	
	// Block until a signal is received
	<-ch
	
	// First we close the connection with MongoDB:
	fmt.Println("\nClosing MongoDB Connection")
	
	// client.Disconnect(context.TODO())	
	if err := client.Disconnect(context.TODO()); err != nil {
        log.Fatalf("Error on disconnection with MongoDB : %v", err)
    }
		
	// Second step : closing the listener
    fmt.Println("Closing the listener")
    if err := lis.Close(); err != nil {
        log.Fatalf("Error on closing the listener : %v", err)
	}
	// Finally, we stop the server
	fmt.Println("Stopping the server")
    s.Stop()
    fmt.Println("End of Program")
}