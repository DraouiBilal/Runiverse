package main

import (
	"context"
	"fmt"
	"log"

	cri "github.com/DraouiBilal/Runiverse-cri/cri/v1"
	"github.com/DraouiBilal/Runiverse/runner/services/container_runtime"
	"github.com/DraouiBilal/Runiverse/runner/services/runner"
)

// Define the server struct
type Server struct {
	// Define any necessary fields
	cri.UnimplementedRuntimeServiceServer
	Runtime container_runtime.ContainerRuntime
}

// Implement the gRPC service methods
func (s *Server) RunCode(ctx context.Context, req *cri.RunCodeRequest) (*cri.RunCodeResponse, error) {
	fmt.Println("Received container creation request:", req)
	container := container_runtime.Container{
		Image:   req.Image,
		Command: req.Command,
		Mounts: []container_runtime.ContainerMount{
			{
				Destination: "/app",
				Source:      "/home/drale/work/open-source/Runiverse/Runiverse-core/services/container-manager/static/",
				Options:     []string{"rbind"},
			},
		},
	}

	logs, err := runner.RunCode(s.Runtime, container)

	if err != nil {
		log.Println(err)
	}

	response := cri.RunCodeResponse{Logs: logs, Err: err.Error()}

	return &response, err
}
