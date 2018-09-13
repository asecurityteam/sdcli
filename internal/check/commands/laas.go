package commands

import (
	"os/exec"
	"strings"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const missingLaas = `
laas is not installed
Logs are centralized into Logging as a Service (LaaS).

In order to access service logs, you need to set up
and sign in to the LaaS CLI:
  https://hello.atlassian.net/wiki/spaces/OBSERVABILITY/pages/140617666/LaaS+CLI

For general information about LaaS, see observability user documentation:
  https://hello.atlassian.net/wiki/spaces/OBSERVABILITY/pages/140622107/User+Documentation`

const unauthenticatedLaas = `
You have laas installed, but you are not signed in.

To sign in, run:
  laas login

For troubleshooting tips, see:
  https://hello.atlassian.net/wiki/spaces/OBSERVABILITY/pages/140617666/LaaS+CLI`

// LaasChecker verifies laas cli is installed, and is authenticated
type LaasChecker struct {
	name string
	r    runner.Runner
}

// NewLaasChecker creates and returns a new LaasChecker
func NewLaasChecker(r runner.Runner) *LaasChecker {
	return &LaasChecker{
		name: "laas",
		r:    r,
	}
}

// Name returns the name of the checker
func (c *LaasChecker) Name() string {
	return c.name
}

// Check executes the checker
// LaasChecker checks to see if the go executable is on the $PATH of the local machine.
// It also verifies that the client is authenticated
func (c *LaasChecker) Check() error {
	var out, err = c.r.Run("laas", "whoami")
	if err != nil {
		// Check to see if laas is installed
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return &CheckerFailure{
				Message: missingLaas,
			}
		}
		return errors.Wrap(err, "error executing laas binary")
	}

	// laas currently doesn't use exit codes or stderr to indicate failure
	if strings.Contains(strings.ToLower(string(out)), "not currently logged in") {
		return &CheckerFailure{
			Message: unauthenticatedLaas,
		}
	}

	return nil
}
