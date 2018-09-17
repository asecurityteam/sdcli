package build

import (
	"bitbucket.org/asecurityteam/sdcli/internal/build/commands"

	"github.com/spf13/cobra"
)

var Verbose bool

// buildCmd represents the build command
func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "build",
		Short: "installs project dependencies, runs tests, and performs Docker builds",
		Long: `build is responsible for enforcing our CI/CD contract. It handles installing
developer dependencies, running tests, generating code coverage, and performing both Docker
builds and deployments.`,
	}

	command.AddCommand(commands.DepCommand())
	command.AddCommand(commands.TestCommand())
	command.AddCommand(commands.DockerCommand())
	return command
}
