package commands

import (
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

// Docker represents common docker commands
type Docker struct {
	r runner.Runner
}

// NewDocker returns new Docker
func NewDocker(r runner.Runner) *Docker {
	return &Docker{r: r}
}

// BuildImage builds an image based off a service
func (d *Docker) BuildImage(service *Service) error {
	out, err := d.r.Run("docker", "build", "-t", service.Image, ".")
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}

// PushImage pushes an image based off a service
func (d *Docker) PushImage(service *Service) error {
	out, err := d.r.Run("docker", "push", service.Image)
	if err != nil {
		return errors.Wrap(err, string(out))
	}
	return nil
}
