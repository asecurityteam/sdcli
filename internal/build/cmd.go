package build

import (
	"os"

	"bitbucket.org/asecurityteam/sdcli/internal/build/commands"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
func NewCommand() *cobra.Command {
	var tag string

	command := &cobra.Command{
		Use:   "build",
		Short: "installs project dependencies, runs tests, and performs Docker builds",
		Long: `build is responsible for enforcing our CI/CD contract. It handles installing
developer dependencies, running tests, generating code coverage, and performing both Docker
builds and deployments.`,
		Run: func(cmd *cobra.Command, args []string) {
			output, err := commands.BuildContainer(tag)
			if err != nil {
				cmd.Printf("Error building Docker image: %s\n", output)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}

	command.Flags().StringVarP(&tag, "tag", "t", "", "Docker image tag")

	command.AddCommand(commands.DepCommand())
	command.AddCommand(commands.TestCommand())
	command.AddCommand(commands.DeployCommand())
	return command
}
