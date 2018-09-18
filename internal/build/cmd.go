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
	command := &cobra.Command{
		Use:   "build",
		Short: "installs project dependencies, runs tests, and performs Docker builds",
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			tag, err := cmd.Flags().GetString("tag")
			if err != nil {
				cmd.Printf("Error getting tag flag: %s", err.Error())
				os.Exit(1)
			}
			output, err := commands.BuildContainer(runner, tag)
			if err != nil {
				cmd.Printf("Error building Docker image: %s\n", output)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	command.Flags().StringP("tag", "t", "", "Docker image tag")

	command.AddCommand(commands.DepCommand())
	command.AddCommand(commands.TestCommand())
	command.AddCommand(commands.NewDeployCommand(r).Command)
	return command
}
