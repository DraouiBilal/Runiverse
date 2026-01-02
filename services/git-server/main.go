package main

import (
	cri "github.com/DraouiBilal/Runiverse-cri/git/v1"
	"github.com/DraouiBilal/Runiverse/git-server/services/git"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {


	var gitServer git.GitServer

	gitServer = &git.Server{RootLocation: "/home/drale/work/trash/git-server"}

	port := os.Getenv("PORT")

	if port == "" {
		port = "50054"
	}


	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	cri.RegisterGitServiceServer(grpcServer, &Server{GitServer: gitServer})

	log.Println("Server is listening on port " + port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
