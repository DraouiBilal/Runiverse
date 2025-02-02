package main

import (
	"log"
	"github.com/DraouiBilal/Runiverse/container_runtime"
	"github.com/DraouiBilal/Runiverse/container_runtime/setup"
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
    runtimes := setup.Setup(true)
    log.Println(runtimes)
	id := runtimes[0].CreateContainer(container_runtime.Container{
		Image:   "golang",
		Command: []string{"go", "run", "/app/main.go"},
		Mounts: []container_runtime.ContainerMount{
			container_runtime.ContainerMount{
				Destination: "/app",
				Source:      "/home/drale/work/open-source/Runiverse/services/container-manager/static",
				Options:     []string{"rbind"},
			},
		},
	})

	id = runtimes[0].StartContainer(container_runtime.Container{Id: id})

    runtimes[0].WaitForContainer(container_runtime.Container{Id: id})

	logs := runtimes[0].GetLogs(container_runtime.Container{Id: id})
	log.Println(logs)
}
