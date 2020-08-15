package main

import (
	"google.golang.org/grpc/codes"
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
	"google.golang.org/grpc/status"
	
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

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error){
	fmt.Printf("Create blog request")
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID: %v", err),
		) 
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
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