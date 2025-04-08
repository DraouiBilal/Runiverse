package main

import (
	"github.com/DraouiBilal/Runiverse/container_runtime/setup"
	"github.com/DraouiBilal/Runiverse-cri/cri"
	"github.com/DraouiBilal/Runiverse/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {

	refreash_runtime := true

	if os.Getenv("REFREASH_RUNTIME") == "false" {
		refreash_runtime = false
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "50051"
	}

	runtimes := setup.Setup(refreash_runtime)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register the gRPC service
	cri.RegisterRuntimeServiceServer(grpcServer, &server.Server{Runtime: runtimes[0]})
	log.Println("Server is listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
