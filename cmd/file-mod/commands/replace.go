package commands

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	findAndReplaceArgs = 3

	findAndReplaceCmd = &cobra.Command{
		Use:   "find-and-replace REGEX_STR STRING FILE",
		Short: "Replace REGEX_STR with STRING in FILE. Supports multiline replace.",
		Args:  cobra.ExactArgs(findAndReplaceArgs),
		RunE:  findAndReplaceE,
	}
)

func init() {
	rootCmd.AddCommand(findAndReplaceCmd)
}

func findAndReplaceE(cmd *cobra.Command, args []string) error {
	regexStr := args[0]
	newString := args[1]
	filePath := args[2]

	r, err := regexp.Compile(regexStr)
	if err != nil {
		return errors.Wrapf(err, "%s is not a valid regular expression", regexStr)
	}

	fileContents, err := fs.ReadFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read contents of file %s", filePath)
	}

	newContents := r.ReplaceAll(fileContents, []byte(newString))

	err = fs.WriteFile(filePath, newContents, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write contents to file %s", filePath)
	}

	return nil
}
