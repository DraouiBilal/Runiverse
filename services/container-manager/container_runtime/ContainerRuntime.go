package container_runtime

type System struct {
	Name       string
	Runtime    string
	SocketPath string
	Version    Version `json:"version"`
}

type Version struct {
	ApiVersion string `json:"APIVersion"`
	Version    string `json:"Version"`
}

type Container struct {
	Id      string           `json:"Id"`
	Name    string           `json:"name"`
	Image   string           `json:"image"`
	Command []string         `json:"command"`
	Mounts  []ContainerMount `json:"mounts"`
}

type ContainerMount struct {
	Source      string   `json:"source"`
	Destination string   `json:"destination"`
	Options     []string `json:"options"`
}

//type Image struct {
//	Id         string
//	name       string
//	repository string
//	tag        string
//	created    string
//	size       string
//}

type ContainerRuntime interface {

	// Runtime
	GetSocket() bool
	GetInfo() System

	// Containers
	CreateContainer(Container) string
	StartContainer(Container) string
	StopContainer(Container) string
	WaitForContainer(Container) string
	GetLogs(Container) string
	GetHealthCheck(Container) string
}
