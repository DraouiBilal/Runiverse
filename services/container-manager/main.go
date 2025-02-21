package main

import (
	"log"
	"github.com/DraouiBilal/Runiverse/container_runtime/setup"
	"github.com/DraouiBilal/Runiverse/cri"
	"github.com/DraouiBilal/Runiverse/server"
	"google.golang.org/grpc"
	"net"
)

func main() {
	runtimes := setup.Setup(false)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register the gRPC service
	cri.RegisterRuntimeServiceServer(grpcServer, &server.Server{Runtime: runtimes[0]})
	log.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

