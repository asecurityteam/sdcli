package check

import (
	"bytes"
	"fmt"
	"regexp"

	"bitbucket.org/asecurityteam/sdcli/internal/check/commands"
	"bitbucket.org/asecurityteam/sdcli/internal/output"
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/spf13/cobra"
)

const long = `
Check runs a series of checks, verifying the proper tools are installed, and that the proper version of those tools are installed.
The following items are checked:
    * git cli is installed along with necessary access to Bitbucket
    * Google-golang at an accepted version
    * Micros CLI (and authentication)
    * Laas CLI (and authentication)
    * Docker is installed at an accepted version
    * Authenticated with docker.atl-paas.net
		* dep is installed at an accepted version
		* golangci-linter is installed`

type multiError struct {
	errors []error
}

func (e *multiError) Error() string {
	var buff = bytes.Buffer{}
	buff.WriteString("Encountered errors while conducting checks:\n")
	for _, err := range e.errors {
		buff.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}
	return buff.String()
}

type checker interface {
	// Check function will check a specific developer dependency
	Check() error
	// Name returns the name of the check
	Name() string
}

// NewCommand returns a new check command
func NewCommand() *cobra.Command {
	var checks []checker
	var r = runner.ExecRunner{}
	checks = append(checks, commands.NewBitbucketChecker(r))
	checks = append(checks, commands.NewGoChecker(r))
	checks = append(checks, commands.NewMicrosChecker(r))
	checks = append(checks, commands.NewLaasChecker(r))
	checks = append(checks, commands.NewDockerChecker(r))
	checks = append(checks, commands.NewRegistryChecker(r))
	checks = append(checks, commands.NewDepChecker(r))
	checks = append(checks, commands.NewLinterChecker(r))

	return &cobra.Command{
		Use:           "check",
		Short:         "Checks to see if the environment of the current machine satisfies SecDev requirements.",
		Long:          long,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          runChecks(checks),
	}
}

func runChecks(checks []checker) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var out = cmd.OutOrStdout()
		var errors []error
		for _, check := range checks {
			fmt.Fprintf(out, "Checking %s... ", check.Name())
			switch err := check.Check(); err.(type) {
			case nil:
				output.Pass(out, "ok")
			case *commands.CheckerFailure:
				var msg = indentBlock(err.Error()+"\n", "    ")
				output.Fail(out, "failure")
				fmt.Fprintln(out, msg)
				errors = append(errors, err)
			default:
				var msg = indentBlock(err.Error()+"\n", "    ")
				output.Fail(out, "error")
				fmt.Fprintln(out, msg)
				errors = append(errors, err)
			}
		}
		if len(errors) == 0 {
			return nil
		}
		return &multiError{errors: errors}
	}
}

func indentBlock(text string, indent string) string {
	rx := regexp.MustCompile(`(?m)(^)(.+)$`)
	return rx.ReplaceAllString(text, indent+"$2")
}
