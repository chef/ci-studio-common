package commands

import (
	"log"
	"os/exec"
	"time"

	"github.com/avast/retry-go"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
)

var (
	archiver filesystem.Archiver

	fs filesystem.FileSystem

	fslock filesystem.Locker

	ciutils install.Install

	rootCmd = &cobra.Command{
		Use:          "ci-utils",
		Short:        "Utility operations to manage the installation of ci-utils",
		SilenceUsage: true,
	}

	execCommand = exec.Command
)

// Execute handles the execution of child commands and flags
func Execute() {
	fs = filesystem.NewOsFs()
	archiver = &filesystem.Unarchiver{}
	ciutils = install.DefaultInstall()

	fslock = &filesystem.OsLock{
		RetryAttempts:  5,
		RetryDelay:     100 * time.Millisecond,
		RetryDelayType: retry.BackOffDelay,
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
