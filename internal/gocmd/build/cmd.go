package build

import (
	"bitbucket.org/asecurityteam/sdcli/internal/gocmd/build/commands"
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const long = `
build is responsible for enforcing our CI/CD contract.
It handles:
	* installing developer dependencies
	* running tests
	* generating code coverage
	* Docker builds and deployments.`

// buildCmd represents the build command
func NewCommand() *cobra.Command {
	r := runner.ExecRunner{}
	d := commands.NewDocker(r)
	command := &cobra.Command{
		Use:   "build",
		Short: "installs project dependencies, runs tests, and performs Docker builds",
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			service, err := commands.NewService(r, true, nil)
			if err != nil {
				return errors.Wrap(err, "error initializing service")
			}
			if err = d.BuildImage(service); err != nil {
				return errors.Wrap(err, "error building image")
			}
			return nil
		},
	}

	command.AddCommand(commands.DepCommand())
	command.AddCommand(commands.NewDeployCommand(r, d).Command)
	return command
}
