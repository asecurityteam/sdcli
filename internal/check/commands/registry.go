package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const missingRegistryTpl = `
You must authenticate docker against our private registry.

The docker containers that we build for SecDev Services
are stored in a private docker registry (%[1]s).

It seems you haven't authenticated with the registry.
To authenticate, run the following command:
          $ docker login %[1]s

You will be prompted for a username and password. Please use
your LDAP/Crowd credentials.`

// RegistryChecker is a checker that verifies that you have access to Atlassian internal docker registries
type RegistryChecker struct {
	name     string
	registry string
	r        runner.Runner
}

// NewRegistryChecker creates and returns a new RegistryChecker
func NewRegistryChecker(r runner.Runner) *RegistryChecker {
	return &RegistryChecker{
		name:     "docker-registry",
		registry: "docker.atl-paas.net",
		r:        r,
	}
}

// Name returns the name for this checker
func (c *RegistryChecker) Name() string {
	return c.name
}

// Check executes the checker
// RegistryChecker checks to see if the system is authenticated to the internal Atlassian
// docker registries
func (c *RegistryChecker) Check() error {
	var failure = &CheckerFailure{Message: fmt.Sprintf(missingRegistryTpl, c.registry)}

	// Open ~/.docker/config.json and we should find an
	// auths["name"].auth key
	var configFile string
	var err error
	if configFile, err = homedir.Expand("~/.docker/config.json"); err != nil {
		return err
	}

	var rawOutput []byte
	rawOutput, err = c.r.Run("cat", configFile)
	switch err.(type) {
	case nil:
	case *exec.ExitError:
		return failure
	default:
		return errors.Wrap(err, "error reading docker config file")
	}

	var out struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
		CredsStore string `json:"credsStore"`
	}

	if err = json.Unmarshal(rawOutput, &out); err != nil {
		return errors.Wrap(err, "error decoding docker config")
	}

	// If there's no auths key for the registry, fail the check
	if _, ok := out.Auths[c.registry]; !ok {
		return failure
	}

	if out.Auths[c.registry].Auth == "" && out.CredsStore == "" {
		return failure
	}

	return nil
}
