// cmd/main.go
package main

import (
	"codespace-clone/internal/dockerops"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func main() {
	fmt.Println("hello codespace")
	var rootCmd = &cobra.Command{Use: "codespace"}
	rootCmd.AddCommand(upCommand())
	rootCmd.AddCommand(downCommand()) // Add the down command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func upCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up [folder-path]",
		Short: "Spin up the development environment",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			folderPath := "."
			if len(args) > 0 {
				folderPath = args[0]
			}
			dockerops.StartCodeServerContainer(folderPath)
		},
	}
}

func downCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop the development environment",
		Run: func(cmd *cobra.Command, args []string) {
			containerName := "codespace-code-server"

			// Stop the container
			fmt.Printf("Stopping the container %s...\n", containerName)
			stopCmd := exec.Command("docker", "stop", containerName)
			if err := stopCmd.Run(); err != nil {
				fmt.Printf("Error stopping container: %v\n", err)
			} else {
				fmt.Printf("Container %s stopped.\n", containerName)
			}

			// Remove the container
			fmt.Printf("Removing the container %s...\n", containerName)
			rmCmd := exec.Command("docker", "rm", containerName)
			if err := rmCmd.Run(); err != nil {
				fmt.Printf("Error removing container: %v\n", err)
			} else {
				fmt.Printf("Container %s removed.\n", containerName)
			}
		},
	}
}
