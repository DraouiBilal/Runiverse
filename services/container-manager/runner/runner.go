package runner

import (
    "github.com/DraouiBilal/Runiverse/container_runtime"
)

func RunCode(runtime container_runtime.ContainerRuntime, container container_runtime.Container) string {

	id := runtime.CreateContainer(container)

	id = runtime.StartContainer(container_runtime.Container{Id: id})

    runtime.WaitForContainer(container_runtime.Container{Id: id})

	logs := runtime.GetLogs(container_runtime.Container{Id: id})

    return logs
}
