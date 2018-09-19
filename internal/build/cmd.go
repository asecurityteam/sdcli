package build

import (
	"os"

	"bitbucket.org/asecurityteam/sdcli/internal/build/commands"
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
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
	var r = runner.ExecRunner{}
	var d = commands.NewDocker(r)
	command := &cobra.Command{
		Use:   "build",
		Short: "installs project dependencies, runs tests, and performs Docker builds",
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			service, err := commands.NewService(r, true, nil)
			if err != nil {
				cmd.Printf("Error initializing service: %s\n", err.Error())
				os.Exit(1)
			}
			if err = d.BuildImage(service); err != nil {
				cmd.Printf("Error building Docker image: %s\n", err.Error())
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	command.AddCommand(commands.DepCommand())
	command.AddCommand(commands.TestCommand())
	command.AddCommand(commands.NewDeployCommand(r, d).Command)
	return command
}
