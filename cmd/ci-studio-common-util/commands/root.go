package commands

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "ci-studio-common-util",
	Short:        "Utility operations to manage the installation of ci-studio-common",
	SilenceUsage: true,
}

// Execute handles the execution of child commands and flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
