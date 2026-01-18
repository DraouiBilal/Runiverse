package runner

import (
	"log"

	"github.com/DraouiBilal/Runiverse/runner/services/container_runtime"
)

func RunCode(runtime container_runtime.ContainerRuntime, container container_runtime.Container) (string, error) {


	id, create_err := runtime.CreateContainer(container)
	log.Println("Container", container, "runtime", runtime)
	
	if create_err != nil {
		return "", create_err
	}

	id, start_err := runtime.StartContainer(container_runtime.Container{Id: id})

	if start_err != nil {
		return "", start_err
	}

	_, wait_err := runtime.WaitForContainer(container_runtime.Container{Id: id})

	if wait_err != nil {
		return "", wait_err
	}

	logs, get_logs_err := runtime.GetLogs(container_runtime.Container{Id: id})

	if get_logs_err != nil {
		return "", get_logs_err
	}

    return logs, nil
}
