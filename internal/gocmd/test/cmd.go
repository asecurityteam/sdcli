package test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	IntegrationFlag = "integration"
	CoverageFlag    = "coverage"
	CoverageFile    = "coverage.out"
)

// NewCommand returns a new test command
func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "test",
		Short: "run unit/integration tests and generate coverage reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			integration, err := cmd.Flags().GetBool(IntegrationFlag)
			if err != nil {
				return errors.Wrap(err, "error getting integration flag")
			}
			coverage, err := cmd.Flags().GetBool(CoverageFlag)
			if err != nil {
				return errors.Wrap(err, "error getting coverage flag")
			}

			testFlags := []string{"test", "-race", "-v", "-cover"}

			if coverage {
				os.Remove(CoverageFile)
				testFlags = []string{"test", "-coverprofile", CoverageFile, "./..."}
				exec.Command("go", testFlags...).CombinedOutput()
				coverageOutput, _ := exec.Command("go", "tool", "cover", fmt.Sprintf("-func=%s", CoverageFile)).CombinedOutput()
				cmd.Printf("%s\n", coverageOutput)
				return nil
			}

			if integration {
				testFlags = append(testFlags, "-tags=integration")
			}

			testFlags = append(testFlags, "./...")
			testOutput, _ := exec.Command("go", testFlags...).CombinedOutput()
			cmd.Printf("%s\n", testOutput)

			return nil
		},
	}

	command.Flags().BoolP(IntegrationFlag, "i", false, "Run integration tests")
	command.Flags().BoolP(CoverageFlag, "c", false, "Display coverage")

	return command
}
