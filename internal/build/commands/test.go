package commands

import (
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	IntegrationFlag = "integration"
	CoverageFlag    = "coverage"
)

// TestCommand returns a new check command
func TestCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "test",
		Short: "run unit/integration tests and generate coverage reports",
		Run: func(cmd *cobra.Command, args []string) {
			integration, err := cmd.Flags().GetBool(IntegrationFlag)
			if err != nil {
				cmd.Printf("Error getting integration flag: %s", err.Error())
			}
			coverage, err := cmd.Flags().GetBool(CoverageFlag)
			if err != nil {
				cmd.Printf("Error getting coverage flag: %s", err.Error())
			}

			testFlags := []string{"test", "-race", "-v", "-cover"}

			if coverage {
				testFlags = []string{"test", "-coverprofile", "cover.out", "./..."}
				exec.Command("go", testFlags...).CombinedOutput()
				coverageOutput, _ := exec.Command("go", "tool", "cover", "-func=coverage.out").CombinedOutput()
				cmd.Printf("%s\n", coverageOutput)
			} else {
				if integration {
					testFlags = append(testFlags, "-tags=integration")
				}
				testFlags = append(testFlags, "./...")
				testOutput, _ := exec.Command("go", testFlags...).CombinedOutput()
				cmd.Printf("%s\n", testOutput)
			}
		},
	}

	command.Flags().BoolP(IntegrationFlag, "i", false, "Run integration tests")
	command.Flags().BoolP(CoverageFlag, "c", false, "Display coverage")

	return command
}
