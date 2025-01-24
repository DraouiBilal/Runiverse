package container_runtime


type Container struct {
	Id      string
	name    string
	Image   string
	Command string
	Created string
	status  int
	Ports   string
}

type Image struct {
	Id         string
	name       string
	repository string
	tag        string
	created    string
	size       string
}

type ContainerRuntime interface {

    // Containers
    CreateContainer(container Container) string
    ListContainers() []Container
    StartContainer(containerID string)
    StopContainer(containerID string)
    RestartContainer(containerID string)
    GetContainer(containerID string)


    // Images
    CreateImage()
    ListImages()
    StartImage()
    StopImage()
    RestartImage()
    GetImage()


    // Volumes
    CreateVolume()
    ListVolumes()
    StartVolume()
    StopVolume()
    RestartVolume()
    GetVolume()
}

