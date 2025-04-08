package server

import (
	"context"
	"fmt"
	"github.com/DraouiBilal/Runiverse/queue"
	"github.com/DraouiBilal/Runiverse-cri/cri"
)

// Define the server struct
type Server struct {
	// Define any necessary fields
	cri.UnimplementedInvokerServiceServer
	Queue *queue.Queue
}

// Implement the gRPC service methods
func (s *Server) Invoke(ctx context.Context, req *cri.InvocationRequest) (*cri.InvocationResponse, error) {
	fmt.Println("Received invokation request:", req)
	s.Queue.AddJob(req)
	return &cri.InvocationResponse{Logs: ""}, nil
}
