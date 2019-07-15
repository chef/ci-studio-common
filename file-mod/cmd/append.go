package cmd

import (
	"github.com/chef/ci-studio-common/lib"

	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append-if-missing STRING FILE",
	Short: "Append STRING to FILE if not already there.",
	Args:  cobra.MinimumNArgs(2),
	Run:   appendIfMissing,
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

func appendIfMissing(cmd *cobra.Command, args []string) {
	stringToAppend := args[0]
	fileToModify := args[1]

	err := lib.AppendIfMissing(fileToModify, stringToAppend)
	lib.Check(err)
}
