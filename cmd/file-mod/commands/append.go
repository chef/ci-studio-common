package commands

import (
	"github.com/spf13/cobra"
)

var appendCmd = &cobra.Command{
	Use:   "append-if-missing STRING FILE",
	Short: "Append STRING to FILE if not already there.",
	Args:  cobra.ExactArgs(2),
	RunE:  appendE,
}

func init() {
	rootCmd.AddCommand(appendCmd)
}

func appendE(cmd *cobra.Command, args []string) error {
	stringToAppend := args[0]
	fileToModify := args[1]

	return fs.AppendIfMissing(fileToModify, []byte(stringToAppend), 0644)
}
