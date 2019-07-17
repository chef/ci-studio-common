package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "file-mod",
	Short:        "Command line utility to modify files.",
	SilenceUsage: true,
}

// Execute handles the execution of child commands and flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
