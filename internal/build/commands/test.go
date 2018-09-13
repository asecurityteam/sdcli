package commands

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var Integration bool
var Coverage bool

// TestCommand returns a new check command
func TestCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "test",
		Short: "run unit/integration tests and generate coverage reports",
		Run: func(cmd *cobra.Command, args []string) {
			testFlags := []string{"test", "-race", "-v", "-cover"}

			if Coverage {
				testFlags = []string{"test", "-coverprofile", "cover.out", "./..."}
				exec.Command("go", testFlags...).CombinedOutput()
				coverageOutput, _ := exec.Command("go", "tool", "cover", "-func=coverage.out").CombinedOutput()
				cmd.Printf("%s\n", coverageOutput)
			} else {
				if Integration {
					testFlags = append(testFlags, "-tags=integration")
				}
				testFlags = append(testFlags, "./...")
				testOutput, _ := exec.Command("go", testFlags...).CombinedOutput()
				cmd.Printf("%s\n", testOutput)
			}
		},
	}

	command.Flags().BoolVarP(&Integration, "integration", "i", false, "Run integration tests")
	command.Flags().BoolVarP(&Coverage, "coverage", "c", false, "Display coverage")

	return command
}
