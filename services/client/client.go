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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RunCodeRequest{
		Image: "golang",
        Command: []string{"go", "run", "/app/main.go"},
	}

	res, err := client.RunCode(ctx, req)
	if err != nil {
		log.Fatalf("Error calling RunCode: %v", err)
	}

	fmt.Printf("Response from server: %s\n", res.Logs)
}

