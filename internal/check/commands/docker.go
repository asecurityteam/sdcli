package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const dockerDaemonFailure = `
The Docker daemon is not reachable.

This could be due to one of two reasons:

* (Linux) Most likely you have not configured the docker command to
  run without sudo access.

  On Linux, the docker command requires sudo access by default.
  Our documentation generally gives examples that invoke docker
  without sudo, since that is what works for Macs.  You can
  configure your system so that sudo is not required for
  specific users.

  To resolve this on Ubuntu, folow the instructions at:
          https://docs.docker.com/engine/installation/linux/ubuntulinux/#/create-a-docker-group

  Depending on how you installed docker, it may be a simple matter
  of adding your user to the docker group, e.g.:
          sudo usermod -aG docker $(id -un)


* The Docker daemon is not running.  You should configure the Docker
  daemon to start automatically.

  On Mac, make sure you're using Docker for Mac and not Docker Machine:
          https://docs.docker.com/docker-for-mac/

  On Ubuntu 16.04, run 'sudo systemctl docker enable' and 'sudo
  systemctl docker start'. For other Linux, follow the
  instructions at:
           https://docs.docker.com/engine/admin/
`

const missingDocker = `
Docker is the container engine we use for services running in micros.

If you are running Linux, you should be able to install docker using
your package manager:

    * sudo apt-get install docker
    * sudo yum install docker

If you are running OSX, you should install Docker for Mac:
    https://docs.docker.com/docker-for-mac/
`

const oldDockerVersionTpl = `
You need to update Docker so that it meets the below minimum required versions:

  - Client version must be at least %s, have %s
  - Server version must be at least %s, have %s
`

// DockerChecker is a checker that verifies that docker is installed and at the correct version
type DockerChecker struct {
	name             string
	minClientVersion string
	minServerVersion string
	r                runner.Runner
}

// NewDockerChecker creates and returns a new DockerChecker
func NewDockerChecker(r runner.Runner) *DockerChecker {
	// Note: semver treats the '-ce' flag on docker versions as "pre-release" versions. If you omit the "-ce" then
	// any comparison against docker versions containing the "-ce" will fail.
	return &DockerChecker{
		name:             "docker",
		minClientVersion: "17.12.0-ce",
		minServerVersion: "17.12.0-ce",
		r:                r,
	}
}

// Name returns the name of the checker
func (c *DockerChecker) Name() string {
	return c.name
}

// Check executes the checker
// Docker checks to see if the go executable is on the $PATH of the local machine.
// It also verifies the system meets the minimum version requirements.
func (c *DockerChecker) Check() error {
	var rawOut, err = c.r.Run("docker", "version", "-f", "{{json .}}")
	switch err := err.(type) {
	case nil:
	case *exec.ExitError:
		return &CheckerFailure{Message: dockerDaemonFailure}
	case *exec.Error:
		if err.Err == exec.ErrNotFound {
			return &CheckerFailure{Message: missingDocker}
		}
		return errors.Wrap(err, "error executing docker binary")
	default:
		return errors.Wrap(err, "error executing docker binary")
	}

	var out struct {
		Client struct {
			Version string
		}
		Server struct {
			Version string
		}
	}

	if err = json.Unmarshal(rawOut, &out); err != nil {
		return errors.Wrap(err, "error decoding version output")
	}

	var serverOk, clientOk bool
	if serverOk, err = checkMinimumVersion(out.Server.Version, c.minServerVersion); err != nil {
		return errors.Wrap(err, "could not check server version")
	}
	if clientOk, err = checkMinimumVersion(out.Client.Version, c.minClientVersion); err != nil {
		return errors.Wrap(err, "could not check client version")
	}

	// If any of the constraints fail, then fail the check
	if !serverOk || !clientOk {
		return &CheckerFailure{
			Message: fmt.Sprintf(oldDockerVersionTpl, c.minClientVersion, out.Client.Version, c.minServerVersion, out.Server.Version),
		}
	}

	return nil
}
