package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "install-buildkite-agent",
	Short:        "Manage the Buildkite Agent installation",
	SilenceUsage: true,
}

// Execute handles the execution of child commands and flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
