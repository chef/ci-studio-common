package commands

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allowCmd)
}

var allowCmd = &cobra.Command{
	Use:   "allow USER",
	Short: "Allow USER to perform certain necessary operations with sudo",
	Args:  cobra.ExactArgs(1),
	RunE:  allowE,
}
