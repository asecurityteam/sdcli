package commands

import (
	"fmt"
	"os"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type DeployCommand struct {
	*cobra.Command
	docker *Docker
	r      runner.Runner
}

func NewDeployCommand(r runner.Runner, docker *Docker) *DeployCommand {
	deployCmd := &DeployCommand{
		docker: docker,
		r:      r,
		Command: &cobra.Command{
			Use:   "deploy",
			Short: "push to docker registry and micros",
		},
	}

	deployCmd.Run = deployCmd.run

	return deployCmd
}

func (d *DeployCommand) run(cmd *cobra.Command, args []string) {
	service, err := NewService(d.r, false, nil)
	if err != nil {
		cmd.Printf("Error initializing service: %s\n", err.Error())
	}

	if err = d.docker.BuildImage(service); err != nil {
		cmd.Printf("Error building image %s: %s\n", service.Image, err.Error())
		os.Exit(1)
	}

	if err = d.Deploy(service); err != nil {
		cmd.Printf("Error deploying %s: %s\n", service.ServiceName, err.Error())
	}

	cmd.Printf("Successfully deployed %s\n", service.ServiceName)
	os.Exit(0)
}

func (d *DeployCommand) Deploy(service *Service) error {
	var err error

	if err = d.docker.BuildImage(service); err != nil {
		return errors.Wrap(err, "error building docker image")
	}

	if err = d.docker.PushImage(service); err != nil {
		return errors.Wrap(err, "error pushing docker image")
	}

	_, err = d.r.RunEnv(
		[]string{fmt.Sprintf("DOCKER_IMAGE=%s", service.ImageName), fmt.Sprintf("DOCKER_TAG=%s", service.ImageTag)},
		"micros",
		"service:deploy",
		service.ServiceName,
		"-f",
		service.ServiceDescriptor,
	)
	if err != nil {
		return errors.Wrap(err, "error deploying to micros\n")
	}

	return nil
}
