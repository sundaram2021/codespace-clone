package dockerops

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func StartCodeServerContainer(folderPath string) {
	ctx := context.Background()
	cli := initializeDockerClient()

	imageName := "codercom/code-server:latest"

	pullImageIfNotExists(cli, ctx, imageName)

	createAndStartContainer(cli, ctx, imageName, folderPath)
}

func initializeDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("Error initializing Docker client:", err)
		os.Exit(1)
	}
	return cli
}

func pullImageIfNotExists(cli *client.Client, ctx context.Context, imageName string) {
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		fmt.Println("Error listing images:", err)
		os.Exit(1)
	}

	imageExists := false
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				imageExists = true
				break
			}
		}
		if imageExists {
			break
		}
	}

	if !imageExists {
		fmt.Printf("Pulling Docker image %s...\n", imageName)
		out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			fmt.Println("Error pulling Docker image:", err)
			os.Exit(1)
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	}
}

func createAndStartContainer(cli *client.Client, ctx context.Context, imageName string, folderPath string) {
	// Verify that the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Printf("The folder path %s does not exist\n", folderPath)
		os.Exit(1)
	}

	// Get the absolute path of the folder
	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		fmt.Printf("Error getting absolute path of %s: %v\n", folderPath, err)
		os.Exit(1)
	}
	folderPath = absFolderPath

	// Remove any existing container with the same name
	containerName := "codespace-code-server"
	_ = cli.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})

	// Create the host configuration with port bindings and volume mounts
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "8080",
				},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: folderPath,
				Target: "/home/coder/project",
			},
		},
	}

	// Create the container configuration
	containerConfig := &container.Config{
		Image: imageName,
		Tty:   true,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
		Env: []string{
			"PASSWORD=", // Disable password authentication
		},
		Cmd: []string{
			"--auth=none",
		},
		WorkingDir: "/home/coder/project",
	}

	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		fmt.Println("Error creating container:", err)
		os.Exit(1)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Println("Error starting container:", err)
		os.Exit(1)
	}

	fmt.Printf("Container %s is started and code-server is accessible at http://localhost:8080\n", resp.ID)
}
