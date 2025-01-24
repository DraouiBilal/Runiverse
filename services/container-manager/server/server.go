package server

import (
	"context"
	"fmt"
    "github.com/DraouiBilal/Runiverse/cri"
)

// Define the server struct
type Server struct {
	// Define any necessary fields
    cri.UnimplementedRuntimeServiceServer
}

// Implement the gRPC service methods
func (s *Server) CreateContainer(ctx context.Context, req *cri.CreateContainerRequest) (*cri.CreateContainerResponse, error) {
	fmt.Println("Received container creation request:", req)
	return &cri.CreateContainerResponse{ContainerId: "SomeID"}, nil
}



