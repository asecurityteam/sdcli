package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func DeployCommand() *cobra.Command {
	var name, serviceDescriptor string

	command := &cobra.Command{
		Use:   "deploy",
		Short: "push to docker registry and micros",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nameSplit := strings.Split(name, ":")
			if len(nameSplit) != 2 {
				cmd.Printf("%s must be formatted as <name>:<tag>\n", name)
			}
			image := nameSplit[0]
			tag := nameSplit[1]

			buildOutput, err := BuildContainer(name)
			if err != nil {
				cmd.Printf("Error building Docker image: %s\n", buildOutput)
				os.Exit(1)
			}

			dockerOutput, err := exec.Command("docker", "push", name).CombinedOutput()
			if err != nil {
				cmd.Printf("Error pushing %s: %s\n", name, dockerOutput)
				os.Exit(1)
			}

			cmd.Printf("Successfully pushed %s\n", name)

			if _, err := os.Stat(serviceDescriptor); os.IsNotExist(err) {
				cmd.Printf("%s does not exist\n", serviceDescriptor)
				os.Exit(1)
			}

			serviceName := args[0]
			cmd.Printf("Preparing to deploy %s", serviceName)
			microsCmd := exec.Command("micros", "service:deploy", serviceName, "-f", serviceDescriptor)
			microsCmd.Env = append(os.Environ(),
				fmt.Sprintf("DOCKER_IMAGE=%s", image),
				fmt.Sprintf("DOCKER_TAG=%s", tag))
			microsOutput, err := microsCmd.CombinedOutput()
			if err != nil {
				cmd.Printf("Error deploying %s to micros: %s\n", serviceName, microsOutput)
				os.Exit(1)
			}

			cmd.Printf("Successfully deployed %s\n", serviceName)
			os.Exit(0)
		},
	}

	command.Flags().StringVarP(&name, "name", "n", "", "image name (required)")
	command.MarkFlagRequired("name")
	command.Flags().StringVarP(&serviceDescriptor, "file", "f", "", "service descriptor file (required)")
	command.MarkFlagRequired("file")

	return command
}

func BuildContainer(tag string) ([]byte, error) {
	dockerBuildCmd := []string{"build"}
	if tag != "" {
		dockerBuildCmd = append(dockerBuildCmd, []string{"-t", tag}...)
	}
	dockerBuildCmd = append(dockerBuildCmd, ".")
	return exec.Command("docker", dockerBuildCmd...).CombinedOutput()
}
