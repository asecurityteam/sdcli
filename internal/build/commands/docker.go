package commands

import (
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

type Docker struct {
	r runner.Runner
}

func NewDocker(r runner.Runner) *Docker {
	return &Docker{r: r}
}

func (d *Docker) BuildImage(service *Service) error {
	out, err := d.r.Run("docker", "build", "-t", service.Image, ".")
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}

func (d *Docker) PushImage(service *Service) error {
	out, err := d.r.Run("docker", "push", service.Image)
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}
