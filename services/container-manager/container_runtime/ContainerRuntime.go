package container_runtime

import (
	"os"
)

type Container struct {
	Id      string           `json:"Id"`
	Name    string           `json: "name"`
	Image   string           `json:"image"`
	Command []string         `json:"command"`
	Mounts  []ContainerMount `json:"mounts"`
}

type ContainerMount struct {
	Source      string      `json:"Source"`
	Destination string      `json:"Destination"`
	BindOptions BindOptions `json:"BindOptions"`
}

type BindOptions struct {
	Propagation string `json:"Propagation"`
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
