package commands

import (
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/files"
)

var appendCmd = &cobra.Command{
	Use:   "append-if-missing STRING FILE",
	Short: "Append STRING to FILE if not already there.",
	Args:  cobra.ExactArgs(2),
	RunE:  appendIfMissing,
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

func appendIfMissing(cmd *cobra.Command, args []string) error {
	stringToAppend := args[0]
	fileToModify := args[1]

	return files.AppendIfMissing(fileToModify, stringToAppend)
}
