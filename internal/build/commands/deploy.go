package commands

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	nameFlag = "name"
	fileFlag = "file"
)

type DeployCommand struct {
	*cobra.Command
	r runner.Runner
}

func NewDeployCommand(r runner.Runner) *DeployCommand {
	deployCmd := &DeployCommand{
		r: r,
		Command: &cobra.Command{
			Use:   "deploy",
			Short: "push to docker registry and micros",
			Args:  cobra.ExactArgs(1),
		},
	}

	deployCmd.Run = deployCmd.run

	deployCmd.Command.Flags().StringP(nameFlag, "n", "", "image name (required)")
	deployCmd.Command.MarkFlagRequired(nameFlag)
	deployCmd.Command.Flags().StringP(fileFlag, "f", "", "service descriptor file (required)")
	deployCmd.Command.MarkFlagRequired(fileFlag)

	return deployCmd
}

func (d *DeployCommand) run(cmd *cobra.Command, args []string) {
	serviceName := args[0]
	name, err := cmd.Flags().GetString(nameFlag)
	if err != nil {
		cmd.Printf("Error getting name flag: %s", err.Error())
		os.Exit(1)
	}
	file, err := cmd.Flags().GetString(fileFlag)
	if err != nil {
		cmd.Printf("Error getting file flag: %s", err.Error())
		os.Exit(1)
	}

	if err = d.Deploy(name, serviceName, file); err != nil {
		cmd.Printf("Error deploying %s: %s", name, err.Error())
		os.Exit(1)
	}

	cmd.Printf("Successfully deployed %s\n", serviceName)
	os.Exit(0)
}

func (d *DeployCommand) Deploy(image, serviceName, serviceDescriptor string) error {
	nameSplit := strings.Split(image, ":")
	if len(nameSplit) != 2 {
		return errors.Errorf("%s must be formatted as <name>:<tag>\n", image)
	}
	name := nameSplit[0]
	tag := nameSplit[1]

	_, err := BuildContainer(d.r, image)
	if err != nil {
		return errors.Wrap(err, "error building Docker image")
	}

	_, err = d.r.Run("docker", "push", image)
	if err != nil {
		return errors.Wrap(err, "error pushing docker image")
	}

	if _, err := os.Stat(serviceDescriptor); os.IsNotExist(err) {
		return errors.Wrap(err, "service descriptor file does not exist\n")
	}

	_, err = d.r.RunEnv(
		[]string{fmt.Sprintf("DOCKER_IMAGE=%s", name), fmt.Sprintf("DOCKER_TAG=%s", tag)},
		"micros",
		"service:deploy",
		serviceName,
		"-f",
		serviceDescriptor,
	)
	if err != nil {
		return errors.Wrap(err, "error deploying to micros\n")
	}
	return nil
}

func BuildContainer(r runner.Runner, image string) ([]byte, error) {
	dockerBuildCmd := []string{"build"}
	if image != "" {
		dockerBuildCmd = append(dockerBuildCmd, []string{"-t", image}...)
	}
	dockerBuildCmd = append(dockerBuildCmd, ".")
	return r.Run("docker", dockerBuildCmd...)
}
