package commands

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// DockerCommand returns a new check command
func DockerCommand() *cobra.Command {
	var Tag string

	command := &cobra.Command{
		Use:   "docker",
		Short: "build and deploy Docker images",
		Run: func(cmd *cobra.Command, args []string) {
			dockerBuildCmd := []string{"build"}
			if Tag != "" {
				dockerBuildCmd = append(dockerBuildCmd, []string{"-t", Tag}...)
			}
			output, err := exec.Command("docker", dockerBuildCmd...).CombinedOutput()
			if err != nil {
				cmd.Printf("Error building Docker image: %s\n", output)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	command.Flags().StringVarP(&Tag, "tag", "t", "", "Docker image tag")

	return command
}
