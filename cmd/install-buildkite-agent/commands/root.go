package commands

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
)

var (
	fs filesystem.FileSystem

	ciutils install.Install

	rootCmd = &cobra.Command{
		Use:          "install-buildkite-agent",
		Short:        "Manage the Buildkite Agent installation",
		SilenceUsage: true,
	}
)

// Execute handles the execution of child commands and flags
func Execute() {
	fs = filesystem.NewOsFs()
	ciutils = install.DefaultInstall()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
