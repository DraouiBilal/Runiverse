package invoker

import (
	"context"
	"fmt"
	"log"
	"time"
    "google.golang.org/grpc/credentials/insecure"
	pb "github.com/DraouiBilal/Runiverse/cri" // Adjust this import to match your project structure
	"google.golang.org/grpc"
)
func invoke(image string, command []string) {
	// Connect to the server
    conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewRuntimeServiceClient(conn)

	// Example: Call CreateContainer method
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.RunCodeRequest{
		Image: image,
        Command: command,
	}

	res, err := client.RunCode(ctx, req)
	if err != nil {
		log.Fatalf("Error calling RunCode: %v", err)
	}

	fmt.Printf("Response from server: %s\n", res.Logs)
}
