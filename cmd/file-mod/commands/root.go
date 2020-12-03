package commands

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
)

var (
	fs filesystem.FileSystem

	rootCmd = &cobra.Command{
		Use:          "file-mod",
		Short:        "Command line utility to modify files.",
		SilenceUsage: true,
	}
)

// Execute handles the execution of child commands and flags.
func Execute() {
	fs = filesystem.NewOsFs()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
