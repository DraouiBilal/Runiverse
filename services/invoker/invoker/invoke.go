package invoker

import (
	"context"
	"log"
	"time"
    "google.golang.org/grpc/credentials/insecure"
	cri "github.com/DraouiBilal/Runiverse-cri/cri/v1" 
	"google.golang.org/grpc"
)
func Invoke(image string, command []string) {
	// Connect to the server
    conn, err := grpc.NewClient("container-manager:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	defer conn.Close()

	client := cri.NewRuntimeServiceClient(conn)

	// Example: Call CreateContainer method
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &cri.RunCodeRequest{
		Image: image,
        Command: command,
	}

	res, err := client.RunCode(ctx, req)
	if err != nil {
		log.Fatalf("Error calling RunCode: %v", err)
	}

	log.Printf("Response from server: %s\n", res.Logs)
}
