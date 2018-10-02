package lint

import (
	"fmt"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/spf13/cobra"
)

// NewCommand returns a new test command
func NewCommand(r runner.Runner) *cobra.Command {
	command := &cobra.Command{
		Use:   "lint",
		Short: "run golangci-lint checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := r.Run("golangci-lint", "run")
			switch err.(type) {
			case nil:
				return nil
			case *exec.ExitError:
				execErr, _ := err.(*exec.ExitError)
				// golangci-lint uses exit code 1 to indicate there was a failed linter check
				// therefore, if exit code is 1, just print the output
				if status := execErr.Sys().(syscall.WaitStatus).ExitStatus(); status == 1 {
					return fmt.Errorf("%s", out)
				}
				return errors.Wrap(err, "error running linter")
			default:
				return errors.Wrap(err, "error running linter")
			}
		},
		SilenceUsage: true,
	}
	return command
}
