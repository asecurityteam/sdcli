package commands

import (
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// DepCommand returns a new check command
func DepCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deps",
		Short: "install developer dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := exec.Command("dep", "ensure").CombinedOutput(); err != nil {
				return errors.Wrap(err, "error ensuring dependcies")
			}
			return nil
		},
	}
}
