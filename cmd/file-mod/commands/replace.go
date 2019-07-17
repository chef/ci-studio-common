package commands

import (
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/files"
)

var replaceCmd = &cobra.Command{
	Use:   "find-and-replace REGEX_STR STRING FILE",
	Short: "Replace REGEX_STR with STRING in FILE. Supports multiline replace.",
	Args:  cobra.ExactArgs(3),
	RunE:  findAndReplace,
}

func init() {
	rootCmd.AddCommand(replaceCmd)
}

func findAndReplace(cmd *cobra.Command, args []string) error {
	regexStr := args[0]
	stringToWrite := args[1]
	fileToModify := args[2]

	return files.FindAndReplace(fileToModify, regexStr, stringToWrite)
}
