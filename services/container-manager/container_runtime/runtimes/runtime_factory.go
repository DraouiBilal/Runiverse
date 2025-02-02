package runtime

import "github.com/DraouiBilal/Runiverse/container_runtime"

func RuntimeFactory(runtime_type string, params container_runtime.System) container_runtime.ContainerRuntime {
    switch runtime_type{
        case "podman":
            return PodmanRuntime{
                RuntimeBase: RuntimeBase{
                    Name: params.Name,   
                    SocketPath: params.SocketPath,
                    Runtime: params.Runtime,
                },
            }
        default:
            return nil
    }
}
