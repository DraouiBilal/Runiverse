package main

import (
	"log"
	"github.com/DraouiBilal/Runiverse/container_runtime"
	"github.com/DraouiBilal/Runiverse/container_runtime/runtime"
	//	   "net"
	//	   "google.golang.org/grpc"
	//		"github.com/DraouiBilal/Runiverse/cri"
	//		"github.com/DraouiBilal/Runiverse/server"
)

//func main() {
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//
//	grpcServer := grpc.NewServer()
//	// Register the gRPC service
//    cri.RegisterRuntimeServiceServer(grpcServer, &server.Server{})
//	log.Println("Server is listening on port 50051")
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}

func main() {
    podman := runtime.PodmanRuntime{}
    podman.SocketPath = "/run/user/1000/podman/podman.sock"
    
    runtimes := []container_runtime.ContainerRuntime{}
    runtimes = append(runtimes, podman)

    id := podman.CreateContainer(container_runtime.Container{Image: "nginx"})
    log.Println(id)
    id = podman.StartContainer(container_runtime.Container{Id: id})
    log.Println(id)
    id = podman.GetLogs(container_runtime.Container{Id: id})
    log.Println(id)
}
