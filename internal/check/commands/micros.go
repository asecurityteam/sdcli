package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const (
	missingMicros = `
micros command not found.  Is it installed?
micros is the Atlassian PaaS we use to manage AWS microservices.

For general information about micros, see:
    https://hello.atlassian.net/wiki/spaces/MICROS/overview

For installation instructions, see:
    https://hello.atlassian.net/wiki/spaces/MICROS/pages/167212704/Micros+CLI`

	oldMicrosVersionTpl = `
Your micros version is %s.  You need to update micros so it is at
least at version %s.

Most likely, all you need to do is run:
    npm install -g @atlassian/micros-cli

But see the confluence page for details:
    https://hello.atlassian.net/wiki/spaces/MICROS/pages/167212704/Micros+CLI`
)

// MicrosChecker is a checker that verifies that the micros cli tool is
// installed and meets the minimum version requirements
type MicrosChecker struct {
	minVersion string
	name       string
	r          runner.Runner
}

// NewMicrosChecker creates and returns a new MicrosChecker
func NewMicrosChecker(r runner.Runner) *MicrosChecker {
	return &MicrosChecker{
		minVersion: "6.1.1",
		name:       "micros",
		r:          r,
	}
}

// Name returns the checker name
func (c *MicrosChecker) Name() string {
	return c.name
}

// Check executes the checker
// MicrosChecker checks to see if the go executable is on the $PATH of the local machine.
// It also verifies the system meets the minimum version requirements.
func (c *MicrosChecker) Check() error {
	var out, err = c.r.Run("micros", "version")
	if err != nil {
		// if the command is missing from $PATH, treat as a failure
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return &CheckerFailure{
				Message: missingMicros,
			}
		}
		return errors.Wrap(err, "error executing micros")
	}
	var ok bool
	var version = strings.TrimSpace(string(out))
	if ok, err = checkMinimumVersion(version, c.minVersion); err != nil {
		return errors.Wrap(err, "could not check micros version")
	}
	if !ok {
		return &CheckerFailure{
			Message: fmt.Sprintf(oldMicrosVersionTpl, version, c.minVersion),
		}
	}
	return nil
}
