package main

import (
	"github.com/DraouiBilal/Runiverse/container_runtime/setup"
	cri "github.com/DraouiBilal/Runiverse-cri/cri/v1"
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

	runtimes, err := setup.Setup(refreash_runtime)

	if err != nil {
		log.Fatal("Error while setting up the runtime: ",err)
	}

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
