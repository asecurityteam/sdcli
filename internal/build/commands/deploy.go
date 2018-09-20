package commands

import (
	"fmt"

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

	deployCmd.RunE = deployCmd.run

	return deployCmd
}

func (d *DeployCommand) run(cmd *cobra.Command, args []string) error {
	service, err := NewService(d.r, false, nil)
	if err != nil {
		return errors.Wrap(err, "error initializing service")
	}

	if err = d.Deploy(service); err != nil {
		return errors.Wrap(err, "error deploying service")
	}

	cmd.Printf("Successfully deployed %s\n", service.ServiceName)
	return nil
}

func (d *DeployCommand) Deploy(service *Service) error {
	if err := d.docker.BuildImage(service); err != nil {
		return errors.Wrap(err, "error building image")
	}

	if err := d.docker.PushImage(service); err != nil {
		return errors.Wrap(err, "error pushing docker image")
	}

	if _, err := d.r.RunEnv(
		[]string{fmt.Sprintf("DOCKER_IMAGE=%s", service.ImageName), fmt.Sprintf("DOCKER_TAG=%s", service.ImageTag)},
		"micros",
		"service:deploy",
		service.ServiceName,
		"-f",
		service.ServiceDescriptor,
	); err != nil {
		return errors.Wrap(err, "error deploying to micros\n")
	}

	return nil
}
