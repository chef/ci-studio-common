package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "install-buildkite-agent",
	Short:        "Manage the Buildkite Agent installation",
	SilenceUsage: true,
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
