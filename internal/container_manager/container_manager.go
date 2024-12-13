package container_manager

import (
	"context"
	"fmt"
	"os/exec"
	"password_generator/internal/types"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func CreateAndStartContainer(opts types.DatabaseOptions) error {
	// Checks if the Docker service is running
	if err := ensureDockerServiceRunning(); err != nil {
		panic(fmt.Sprintf("Failed to start docker service: %v", err))
	}

	// Creating new API client
	apiClient, err := newAPIClient()
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	//Checks if such a container exists already and if it is running or not
	containerExists, containerIsRunning := checkContainerStatus(apiClient, opts.ContainerName)

	if !containerExists {
		createContainer(apiClient, opts)
	} else {
		if !containerIsRunning {
			startContainer(apiClient, opts.ContainerName)
		} else if containerExists {
			//NOTE fmt.Println("Container '", opts.ContainerName, "' is already running")
		} else {
			fmt.Println("Container '", opts.ContainerName, "' not found. Please create the container first")
		}
	}
	return nil
}

func ensureDockerServiceRunning() error {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "docker")
	if err := cmd.Run(); err == nil {
		return nil
	}

	cmd = exec.Command("sudo", "systemctl", "start", "docker")
	return cmd.Run()
}

func newAPIClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv)
}

func checkContainerStatus(apiClient *client.Client, containerName string) (bool, bool) {
	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(fmt.Sprintf("Failed to get containers list: %v", err))
	}

	var containerExists, containerIsRunning bool

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				containerExists = true
				containerIsRunning = container.State == "running"
				break
			}
		}
		if containerExists {
			break
		}
	}
	return containerExists, containerIsRunning
}

func createContainer(apiClient *client.Client, options types.DatabaseOptions) {
	//NOTE fmt.Println("Creating a container...")

	// Exposing postgres port for connecting to it from outside the container
	exposedPort := nat.Port("5432" + "/tcp")

	resp, err := apiClient.ContainerCreate(context.Background(), &container.Config{
		Image: "postgres",
		Cmd: []string{
			"postgres",
			"-c", "password_encryption=scram-sha-256",
			"-c", "listen_addresses=*",
			"-c", "hba_file=/var/lib/postgresql/data/pg_hba.conf",
		},
		Env: []string{
			"POSTGRES_USER=" + options.UserName,
			"POSTGRES_PASSWORD=" + options.DbPassword,
		},
		ExposedPorts: nat.PortSet{
			exposedPort: struct{}{},
		},
	}, &container.HostConfig{
		//HostPort --> 5432 port
		PortBindings: nat.PortMap{
			exposedPort: []nat.PortBinding{
				{HostPort: options.HostPort},
			},
		},
	}, nil, nil, options.ContainerName)

	if err != nil {
		panic(err)
	}

	if err := apiClient.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
		panic(fmt.Sprintf("Failed to start a container: %v", err))
	}

	//NOTE fmt.Println("Container", options.ContainerName, "is ready to serve")
}

func startContainer(apiClient *client.Client, containerName string) {
	fmt.Println("Container '", containerName, "' exists but is not running. Starting it...")

	if err := apiClient.ContainerStart(context.Background(), containerName, container.StartOptions{}); err != nil {
		panic(fmt.Sprintf("Failed to start a container: %v", err))
	}

	fmt.Println("Container started successfully")
}
