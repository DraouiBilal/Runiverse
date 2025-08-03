package runtime

import (
	"github.com/DraouiBilal/Runiverse/container_runtime"
	"github.com/DraouiBilal/Runiverse-backend-lib/api"
)

type PodmanRuntime struct {
	RuntimeBase
}

func (p PodmanRuntime) GetSocket() (bool, error) {
	return p.SocketExists()
}

func (p PodmanRuntime) CreateContainer(container container_runtime.Container) (string, error) {
	res, err := api.Post[container_runtime.Container]("http://localhost/v5.0.0/libpod/containers/create", container, api.Options{Socket: p.SocketPath})
	return res.Id, err
}

func (p PodmanRuntime) StartContainer(container container_runtime.Container) (string, error) {
	_, err := api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/start", container, api.Options{Socket: p.SocketPath})
	return container.Id, err
}

func (p PodmanRuntime) StopContainer(container container_runtime.Container) (string, error) {
	_, err := api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/stop", container, api.Options{Socket: p.SocketPath})
	return container.Id, err
}

func (p PodmanRuntime) WaitForContainer(container container_runtime.Container) (string, error) {
	_, err := api.Post[interface{}]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/wait", container, api.Options{Socket: p.SocketPath})
    return container.Id + " Done", err
}

func (p PodmanRuntime) GetLogs(container container_runtime.Container) (string, error) {
	logs, err := api.Get[string]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/logs?stdout=true&stderr=true", container, api.Options{Socket: p.SocketPath, Follow: true})
	return *logs, err
}

func (p PodmanRuntime) GetInfo() (container_runtime.System, error) {
	info, err := api.Get[container_runtime.System]("http://localhost/v5.0.0/libpod/info", nil, api.Options{Socket: p.SocketPath})
    info.Name = p.Name
    info.Runtime = p.Runtime
    info.SocketPath = p.SocketPath
    return *info, err
}

func (p PodmanRuntime) GetHealthCheck(container container_runtime.Container) (string, error) {
	logs, err := api.Get[string]("http://localhost/v5.0.0/libpod/containers/"+container.Id+"/logs?stdout=true&stderr=true", container, api.Options{Socket: p.SocketPath, Follow: true})
	return *logs, err
}
