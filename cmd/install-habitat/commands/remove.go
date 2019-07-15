package commands

import (
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Completely uninstall Chef Habitat from the system.",
	RunE:  uninstallHabitat,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
