package cri

import (
	"context"
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
//	"github.com/DraouiBilal/Runiverse  // Import the generated package
)

// Define the server struct
type server struct {
	// Define any necessary fields
    UnimplementedRuntimeServiceServer
}

// Implement the gRPC service methods
func (s *server) CreateContainer(ctx context.Context, req *CreateContainerRequest) (*CreateContainerResponse, error) {
	// Your container creation logic here
	fmt.Println("Received container creation request:", req)
	return &CreateContainerResponse{ContainerId: "Created"}, nil
}

// Start the gRPC server
func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register the gRPC service
    RegisterRuntimeServiceServer(grpcServer, &server{})
	fmt.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


