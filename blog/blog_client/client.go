package main

import (
	"context"
	"fmt"
	//"io"
	"log"
	//"time"

	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/credentials"

	"../blogpb"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:5000", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	
	fmt.Printf("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Stephane",
		Title: "My First Blog",
		Content: "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
}