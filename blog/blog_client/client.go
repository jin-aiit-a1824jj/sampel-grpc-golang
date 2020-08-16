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
	
	// create Blog
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
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	fmt.Println("Reading the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5f38b85033956088d5809ee8a"})
	if err2 != nil {
		fmt.Printf("Error happended while reading: %v \n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happend while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)
}