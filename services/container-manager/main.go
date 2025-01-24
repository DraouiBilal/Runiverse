package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
	"github.com/DraouiBilal/Runiverse/cri"
	"github.com/DraouiBilal/Runiverse/server"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register the gRPC service
    cri.RegisterRuntimeServiceServer(grpcServer, &server.Server{})
	log.Println("Server is listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
