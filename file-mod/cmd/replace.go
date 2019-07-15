package cmd

import (
	"github.com/chef/ci-studio-common/lib"

	"github.com/spf13/cobra"
)

var replaceCmd = &cobra.Command{
	Use:   "find-and-replace REGEX_STR STRING FILE",
	Short: "Replace REGEX_STR with STRING in FILE. Supports multiline replace.",
	Args:  cobra.MinimumNArgs(3),
	Run:   findAndReplace,
}

func init() {
	rootCmd.AddCommand(replaceCmd)
}

func findAndReplace(cmd *cobra.Command, args []string) {
	regexStr := args[0]
	stringToWrite := args[1]
	fileToModify := args[2]

	err := lib.FindAndReplace(fileToModify, regexStr, stringToWrite)
	lib.Check(err)
}
