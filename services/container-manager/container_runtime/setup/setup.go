package setup

import (
	"fmt"
	"github.com/DraouiBilal/Runiverse/container_runtime"
	runtime "github.com/DraouiBilal/Runiverse/container_runtime/runtimes"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
)

func Setup(refreash bool) []container_runtime.ContainerRuntime{
	if !refreash {
        log.Println("Loading config from" + RUNTIME_CONFIG_PATH + RUNTIME_CONFIG_FILE)
        runtimes := []container_runtime.ContainerRuntime{}
        config := loadConfigFile()
        for _, cf := range config {
            runtimes = append(runtimes, runtime.RuntimeFactory(cf.Runtime, container_runtime.System{
                Name: cf.Name,
                Runtime: cf.Runtime,
                SocketPath: cf.SocketPath,
                Version: container_runtime.Version{
                    Version: cf.Version,
                    ApiVersion: cf.APIVersion,
                },
            }))
        }
        return runtimes
	}
	log.Println("Setting Up Runiverse...")

	checkAndCreateRuniverseDir(RUNTIME_CONFIG_PATH)

	runtimes := checkContainerRuntime()

	if len(runtimes) == 0 {
		log.Println("No default runtime found, Checking runtime.yml file for custom config")
	}

	log.Print("Found " + strconv.Itoa(len(runtimes)) + " Runtime")

	for _, rt := range runtimes {
		info := rt.GetInfo()
		fmt.Println(fmt.Sprintf(`
            Name: %s
            Runtime: %s
            Socket path: %s
            Version: %s
            APIVersion: %s
        `, info.Name, info.Runtime, info.SocketPath, info.Version.Version, info.Version.ApiVersion))
	}

	log.Println("Creating runtime config file at " + RUNTIME_CONFIG_PATH + RUNTIME_CONFIG_FILE)
	createConfigFile(runtimes)
    return runtimes
}

func checkContainerRuntime() []container_runtime.ContainerRuntime {
	runtimes := []container_runtime.ContainerRuntime{}

	podman := runtime.PodmanRuntime{
		RuntimeBase: runtime.RuntimeBase{
			Name:       "Podman",
			Runtime:    "podman",
			SocketPath: os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock",
		},
	}

	if podman.SocketExists() {
		runtimes = append(runtimes, podman)
	}

	return runtimes

}

func checkAndCreateRuniverseDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Folder does not exist, create it
		err := os.MkdirAll(path, 0755) // Creates parent directories if needed
		if err != nil {
			return err
		}
		log.Println("Runiverse directory created:", path)
	} else {
		log.Println("Runiverse directory already exists:", path)
	}
	return nil
}

func createConfigFile(runtimes []container_runtime.ContainerRuntime) {
	config := []Config{}

	for _, rt := range runtimes {
		info := rt.GetInfo()
		cf := Config{
			Name:       info.Name,
			Runtime:    info.Runtime,
			SocketPath: info.SocketPath,
			Version:    info.Version.Version,
			APIVersion: info.Version.ApiVersion,
		}
		config = append(config, cf)
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return
	}

	// Create and write to the file
	filePath := RUNTIME_CONFIG_PATH + RUNTIME_CONFIG_FILE
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	log.Println("Runtime config created successfully")
}

func loadConfigFile() []Config{
    data, err := os.ReadFile(RUNTIME_CONFIG_PATH + RUNTIME_CONFIG_FILE)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// Parse the YAML data into a struct
	var config []Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error unmarshaling YAML:", err)
	}

    return config
}

type Config struct {
	Name       string `yaml:"name"`
	Runtime    string `yaml:"runtime"`
	SocketPath string `yaml:"socket_path"`
	Version    string `yaml:"version"`
	APIVersion string `yaml:"api_version"`
}
