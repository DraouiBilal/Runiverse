package main

import (
	"github.com/DraouiBilal/Runiverse-cri/cri"
	"github.com/DraouiBilal/Runiverse/queue"
	"github.com/DraouiBilal/Runiverse/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	q := queue.InitQueue()

	//q.AddJob(&cri.InvocationRequest{
	//	Image:   "golang",
	//	Command: []string{"go", "run", "/app/main.go"},
	//})

	port := os.Getenv("PORT")

	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register the gRPC service
	cri.RegisterInvokerServiceServer(grpcServer, &server.Server{Queue: q})

	log.Println("Server is listening on port " + port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	q.StopWorkers()

	log.Println("Queue closed, all jobs completed")
}
