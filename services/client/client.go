package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/DraouiBilal/Runiverse/cri" // Adjust this import to match your project structure
	"google.golang.org/grpc"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Use WithTransportCredentials for production
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRuntimeServiceClient(conn)

	// Example: Call CreateContainer method
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateContainerRequest{
		ContainerId: "test-container",
	}

	res, err := client.CreateContainer(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CreateContainer: %v", err)
	}

	fmt.Printf("Response from server: %s\n", res.ContainerId)
}

