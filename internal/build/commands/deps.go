package commands

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// DepCommand returns a new check command
func DepCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "deps",
		Short: "install developer dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			output, err := exec.Command("dep", "ensure").CombinedOutput()
			if err != nil {
				cmd.Printf("Error ensuring dependencies: %s\n", output)
				os.Exit(1)
			}
			os.Exit(0)
		},
	}
}
