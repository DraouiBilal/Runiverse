package runtime

import (
	"github.com/DraouiBilal/Runiverse/container_runtime"
	"github.com/DraouiBilal/Runiverse/lib/api"
)

type PodmanRuntime struct {
	RuntimeBase
}

func (p PodmanRuntime) GetSocket() bool {
	return p.SocketExists()
}

func (p PodmanRuntime) CreateContainer(container container_runtime.Container) string {
	res := api.Post[container_runtime.Container]("http://localhost/v5.0.0/libpod/containers/create", container, api.Options{Socket: p.SocketPath})
	return res.Id
}

func (p PodmanRuntime) StartContainer(container container_runtime.Container) string {
    api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/start", container, api.Options{Socket: p.SocketPath})
	return container.Id
}

func (p PodmanRuntime) StopContainer(container container_runtime.Container) string {
	api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/stop", container, api.Options{Socket: p.SocketPath})
	return container.Id
}

func (p PodmanRuntime) WaitForContainer(container container_runtime.Container) string {
    api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/wait", container, api.Options{Socket: p.SocketPath})
    return container.Id + " Done"
}

func (p PodmanRuntime) GetLogs(container container_runtime.Container) string {
	logs := api.Get[string]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/logs?stdout=true&stderr=true", container, api.Options{Socket: p.SocketPath, Follow: true})
	return *logs
}

func (p PodmanRuntime) GetInfo() container_runtime.System {
	info := api.Get[container_runtime.System]("http://localhost/v5.0.0/libpod/info", nil, api.Options{Socket: p.SocketPath})
    info.Name = p.Name
    info.Runtime = p.Runtime
    info.SocketPath = p.SocketPath
    return *info
}

func (p PodmanRuntime) GetHealthCheck(container container_runtime.Container) string {
	logs := api.Get[string]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/logs?stdout=true&stderr=true", container, api.Options{Socket: p.SocketPath, Follow: true})
	return *logs
}
