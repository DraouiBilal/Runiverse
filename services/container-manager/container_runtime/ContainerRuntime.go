package container_runtime

import (
	"os"
)

type Container struct {
	Id      string `json:"Id"`
	Name    string `json: "Name"`
	Image   string `json:"Image"`
	//Command string `json:"Command"`
	//Created string `json:"Created"`
	//Status  int    `json:"Status"`
	//Ports   string `json:"Ports"`
}

type Image struct {
	Id         string
	name       string
	repository string
	tag        string
	created    string
	size       string
}

type RuntimeBase struct {
	SocketPath string
}

func (runtime RuntimeBase) SocketExists() bool {
	_, err := os.Stat(runtime.SocketPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type ContainerRuntime interface {

	// Runtime
	GetSocket() bool

	// Containers
	CreateContainer(Container) string
	// ListContainers() []Container
	StartContainer(Container) string
    StopContainer(Container) string
    GetLogs(Container) string
	// RestartContainer(containerID string)
	// GetContainer(containerID string)

	// Images
	// CreateImage()
	// ListImages()
	// StartImage()
	// StopImage()
	// RestartImage()
	// GetImage()

	// Volumes
	// CreateVolume()
	// ListVolumes()
	// StartVolume()
	// StopVolume()
	// RestartVolume()
	// GetVolume()
}
