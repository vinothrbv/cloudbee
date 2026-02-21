// Small client to test CreatePost (and GetPost). Run the server first, then: go run ./cmd/client
package main

import (
	"context"
	"log"
	"time"

	"crud/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("dial: %v", err)
	}

	defer conn.Close()

	client := pb.NewPostServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a post (id=0 lets server assign)
	now := timestamppb.Now()
	createResp, err := client.CreatePost(ctx, &pb.CreatePostRequest{
		Id:              0,
		Title:           "First post",
		Content:         "Hello from gRPC client.",
		Author:          "alice",
		Tags:            []string{"go", "grpc"},
		PublicationDate: now,
	})
	if err != nil {
		log.Fatalf("CreatePost: %v", err)
	}
	log.Printf("CreatePost OK: id=%d title=%q", createResp.GetPost().GetId(), createResp.GetPost().GetTitle())

	// Create a post specifying the id
	createResp, err = client.CreatePost(ctx, &pb.CreatePostRequest{
		Id:              2,
		Title:           "Second post",
		Content:         "Hello again",
		Author:          "trial",
		Tags:            []string{"grpc", "trial2"},
		PublicationDate: now,
	})
	if err != nil {
		log.Fatalf("CreatePost: %v", err)
	}
	log.Printf("CreatePost OK: id=%d title=%q", createResp.GetPost().GetId(), createResp.GetPost().GetTitle())

	// Fetch it back
	id := createResp.GetPost().GetId()
	getResp, err := client.GetPost(ctx, &pb.GetPostRequest{Id: id})
	if err != nil {
		log.Fatalf("GetPost: %v", err)
	}
	log.Printf("GetPost OK: id=%d title=%q author=%q", getResp.GetPost().GetId(), getResp.GetPost().GetTitle(), getResp.GetPost().GetAuthor())

	// Update the tags to second post
	updateResp, err := client.UpdatePost(ctx, &pb.UpdatePostRequest{
		Id:   id,
		Tags: []string{"go", "grpc", "update"},
	})
	if err != nil {
		log.Fatalf("UpdatePost: %v", err)
	}
	log.Printf("UpdatePost OK: id=%d tags=%v", updateResp.GetPost().GetId(), updateResp.GetPost().GetTags())

	// Delete the first post by specifying the id 1
	deleteResp, err := client.DeletePost(ctx, &pb.DeletePostRequest{Id: 1})
	if err != nil {
		log.Fatalf("DeletePost: %v", err)
	}
	log.Printf("DeletePost OK: deleted = %t", deleteResp.Deleted)
}
