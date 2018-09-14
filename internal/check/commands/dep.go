package commands

import (
	"fmt"
	"os/exec"
	"regexp"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const missingDep = `
dep command not found. Is it on the $PATH?",
We use dep to manage our Go dependencies.

Dep helps us have reproducible builds.

You can get Dep via:
          https://golang.github.io/dep/docs/installation.html
`

const unexpectedDepOutputTpl = `
Running 'dep version' returned '%s', which doesn't seem right.

Are you sure you have Dep installed properly?`

const oldDepVerstionTpl = `
You must have a dep version of at least %s, but you have %s.

You can get Dep via:
          https://golang.github.io/dep/docs/installation.html`

// DepChecker is a checker that verifies that dep is installed and at a minimum version
type DepChecker struct {
	name       string
	minVersion string
	r          runner.Runner
}

// NewDepChecker creates and returns a new dep checker
func NewDepChecker(r runner.Runner) *DepChecker {
	return &DepChecker{
		name:       "dep",
		minVersion: "0.5.0",
		r:          r,
	}
}

// Name returns the name of the checker
func (c *DepChecker) Name() string {
	return c.name
}

// dep outputs its version as:
// dep:
//  version     : v0.5.0
//  build date  : 2018-07-26
//  git hash    : 224a564
//  go version  : go1.10.3
//  go compiler : gc
//  platform    : darwin/amd64
//  features    : ImportDuringSolve=false
var depVersionRx = regexp.MustCompile(`version\s*:\s*v?(\d\.\d(?:\.?\d)?)\b`)

// Check executes the checker
// DepChecker checks to see if the dep executable is on the $PATH of the local machine.
// It also verifies the system meets the minimum version requirements.
func (c *DepChecker) Check() error {
	var out, err = c.r.Run("dep", "version")
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return &CheckerFailure{Message: missingDep}
		}
		return errors.Wrap(err, "error executing dep")
	}

	var matches = depVersionRx.FindStringSubmatch(string(out))
	if matches == nil || len(matches) != 2 {
		return &CheckerFailure{Message: fmt.Sprintf(unexpectedDepOutputTpl, out)}
	}

	var ok bool
	if ok, err = checkMinimumVersion(matches[1], c.minVersion); err != nil {
		return errors.Wrap(err, "could not check dep version")
	}

	if !ok {
		return &CheckerFailure{
			Message: fmt.Sprintf(oldDepVerstionTpl, c.minVersion, matches[1]),
		}
	}

	return nil
}
