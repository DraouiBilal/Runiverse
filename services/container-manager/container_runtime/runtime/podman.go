package runtime

import (
	"log"

	"github.com/DraouiBilal/Runiverse/container_runtime"
	"github.com/DraouiBilal/Runiverse/lib/api"
)

type PodmanRuntime struct {
    container_runtime.RuntimeBase
}

func (p PodmanRuntime) GetSocket() bool{
    return p.SocketExists()
}

func (p PodmanRuntime) CreateContainer(container container_runtime.Container) string {
    res := api.Post[container_runtime.Container]("http://localhost/v4.0.0/libpod/containers/create", container, api.Options{Socket: p.SocketPath})
    log.Println(res)
    return res.Id
}

func (p PodmanRuntime) StartContainer(container container_runtime.Container) string {
    api.Post[interface{}]("http://localhost/v4.0.0/libpod/containers/"+container.Id+"/start", container, api.Options{Socket: p.SocketPath})
    return container.Id
}

func (p PodmanRuntime) StopContainer(container container_runtime.Container) string {
    api.Post[interface{}]("http://localhost/v4.0.0/libpod/containers/"+container.Id+"/stop", container, api.Options{Socket: p.SocketPath})
    return container.Id
}

func (p PodmanRuntime) GetLogs(container container_runtime.Container) string {
    logs := api.Get[string]("http://localhost/v4.0.0/libpod/containers/"+container.Id+"/logs?stdout=true&stderr=true", container, api.Options{Socket: p.SocketPath})
    return *logs
}
