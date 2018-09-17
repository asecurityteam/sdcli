package commands

import (
	"fmt"
	"os/exec"
	"regexp"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const (
	missingGoTpl = `
You'll need to install at least Go version %s. You can get Go
through your package manager, or get a prebuilt distribution
for your platform from https://golang.org/dl`

	unexpectedGoOutputTpl = `
Running 'go version' returned '%s', which doesn't seem right.
Are you sure you have Go installed properly?`

	oldGoVersionTpl = `
SecDev Services targets Go %s or greater, but you have %s.
You can update Go by downloading the latest release at https://golang.org/dl`
)

// GoChecker is a checker that verifies that Go is installed and at a minimum version
type GoChecker struct {
	minVersion string
	name       string
	r          runner.Runner
}

// NewGoChecker creates and returns a new checker
func NewGoChecker(r runner.Runner) *GoChecker {
	return &GoChecker{
		minVersion: "1.11",
		name:       "go",
		r:          r,
	}
}

// Name returns the checker name
func (c *GoChecker) Name() string {
	return c.name
}

var goVersionRx = regexp.MustCompile(`^go version (?:go)?(\d\.\d(?:\.?\d)?)\b`)

// Check executes the checker
// GoChecker checks to see if the go executable is on the $PATH of the local machine.
// It also verifies the system meets the minimum version requirements.
func (c *GoChecker) Check() error {
	var out, err = c.r.Run("go", "version")
	if err != nil {
		// If the returned error is an exec.Error whose Err property is an exec.ErrNotFound, then we can turn that into a failed check
		if e, ok := err.(*exec.Error); ok && e.Err == exec.ErrNotFound {
			return &CheckerFailure{
				Message: fmt.Sprintf(missingGoTpl, c.minVersion),
			}
		}
		// Otherwise, it's a general unknown failure so we can wrap it with some more context
		return errors.Wrap(err, "error executing go binary")
	}

	// go version should print out `go version goX.Y`
	var matches = goVersionRx.FindStringSubmatch(string(out))
	if matches == nil || len(matches) != 2 {
		return &CheckerFailure{
			Message: fmt.Sprintf(unexpectedGoOutputTpl, string(out)),
		}
	}

	var ok bool
	if ok, err = checkMinimumVersion(matches[1], c.minVersion); err != nil {
		return errors.Wrap(err, "could not check go version")
	}

	if !ok {
		return &CheckerFailure{
			Message: fmt.Sprintf(oldGoVersionTpl, c.minVersion, matches[1]),
		}
	}

	return nil
}
