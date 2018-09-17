package commands

import (
	"os/exec"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const lintMissing = `
golangci-lint command not found. Is it on the $PATH?",
SecDev uses golangci-lint to handle many code quality checks.

To install golangci-lint, you must have Go already set up, and then:
          go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

If you've set up Go as recommended, your PATH will include $GOPATH/bin
so you can use golangci-lint directly without specifying the full path.`

// LinterChecker checks to see if golangci-lint is installed on the $PATH
type LinterChecker struct {
	name string
	r    runner.Runner
}

// NewLinterChecker creates and returns a new LinterChecker
func NewLinterChecker(r runner.Runner) *LinterChecker {
	return &LinterChecker{
		name: "golangci-lint",
		r:    r,
	}
}

// Name returns the name of the checker
func (c *LinterChecker) Name() string {
	return c.name
}

// Check executes the checker
// GoChecker checks to see if the go executable is on the $PATH of the local machine.
// It also verifies the system meets the minimum version requirements.
func (c *LinterChecker) Check() error {
	var _, err = c.r.Run("golangci-lint", "linters")
	if err != nil {
		if e, ok := err.(*exec.Error); ok && e.Err == exec.ErrNotFound {
			return &CheckerFailure{Message: lintMissing}
		}
		// Otherwise, it's a general unknown failure so we can wrap it with some more context
		return errors.Wrap(err, "error executing gometalinter binary")
	}

	return nil
}
