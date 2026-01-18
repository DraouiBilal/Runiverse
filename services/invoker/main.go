package main

import (
	cri "github.com/DraouiBilal/Runiverse-cri/cri/v1"
	"github.com/DraouiBilal/Runiverse/invoker/services/queue"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	q := queue.InitQueue()

	port := os.Getenv("PORT")

	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	cri.RegisterInvokerServiceServer(grpcServer, &Server{Queue: q})

	log.Println("Server is listening on port " + port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	q.StopWorkers()

	log.Println("Queue closed, all jobs completed")
}
