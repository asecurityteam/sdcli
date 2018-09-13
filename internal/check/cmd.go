package check

import (
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
    * Authenticated with docker.atl-paas.net`

type checker interface {
	// Check function will check a specific developer dependency
	Check() error
	// Name returns the name of the check
	Name() string
}

// Command is a struct representing the command which executes all of the development environment checks
type Command struct {
	*cobra.Command
	checks []checker
}

// NewCommand returns a new check command
func NewCommand() *Command {
	var cmd = &Command{
		Command: &cobra.Command{
			Use:   "check",
			Short: "Checks to see if the environment of the current machine satisfies SecDev requirements.",
			Long:  long,
		},
	}
	cmd.Run = cmd.run

	// register new checkers here
	var r = runner.ExecRunner{}
	cmd.checks = append(cmd.checks, commands.NewBitbucketChecker(r))
	cmd.checks = append(cmd.checks, commands.NewGoChecker(r))
	cmd.checks = append(cmd.checks, commands.NewMicrosChecker(r))

	return cmd
}

func (c *Command) run(cmd *cobra.Command, args []string) {
	var out = c.OutOrStdout()
	for _, check := range c.checks {
		fmt.Fprintf(out, "Checking %s... ", check.Name())
		switch err := check.Check(); err.(type) {
		case nil:
			output.Check(out, "ok")
		case *commands.CheckerFailure:
			var msg = indentBlock(err.Error()+"\n", "    ")
			output.Fail(out, "failure")
			fmt.Println(msg)
		default:
			var msg = indentBlock(err.Error()+"\n", "    ")
			output.Fail(out, "error")
			fmt.Fprintln(out, msg)
		}
	}
}

func indentBlock(text string, indent string) string {
	rx := regexp.MustCompile(`(?m)(^)(.+)$`)
	return rx.ReplaceAllString(text, indent+"$2")
}
