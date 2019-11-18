package commands

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zbiljic/go-filelock"
	"go.uber.org/multierr"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
)

type rootCmdOptions struct {
	channel string
	version string
}

var (
	ciutils install.Install

	fs filesystem.FileSystem

	fslock filesystem.Locker

	execCommand = exec.Command

	// DO NOT MODIFY -- This value is automatically updated by Expeditor
	globalHabValue = "0.90.6"

	rootCmd = &cobra.Command{
		Use:   "install-habitat",
		Short: "Install Chef Habitat",
		RunE:  rootE,
	}

	rootCmdOpts = &rootCmdOptions{}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&rootCmdOpts.version, "version", "v", globalHabValue, "Which version of Habitat you wish to install.")
	rootCmd.Flags().StringVarP(&rootCmdOpts.channel, "channel", "c", "stable", "The channel from which you wish to install Habitat.")
	rootCmd.Flags().StringP("target", "t", "", "The kernel target for this installation.")
	if err := viper.BindPFlag("target", rootCmd.Flags().Lookup("target")); err != nil {
		panic(err)
	}
	viper.SetDefault("target", fmt.Sprintf("x86_64-%s", runtime.GOOS))
}

// Execute handles the execution of child commands and flags
func Execute() {
	fs = filesystem.NewOsFs()
	ciutils = install.DefaultInstall()
	fslock = &filesystem.OsLock{}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootE(cmd *cobra.Command, args []string) error {
	lockfile := ciutils.LockPath("install-habitat")
	lock, err := fslock.GetLock(lockfile)
	if err != nil {
		return errors.Wrapf(err, "could not create %s", lockfile)
	}

	if _, err = lock.TryLock(); err == nil {
		return maybeInstallHabitat(cmd, lock)
	} else if err == filesystem.ErrLocked {
		cmd.Println("Chef Habitat install already in progress -- skipping")
	} else {
		return errors.Wrap(err, "could not acquire file lock")
	}

	return nil
}

func maybeInstallHabitat(cmd *cobra.Command, lock filelock.TryLockerSafe) error {
	currentVersion := "none"

	output, err := execCommand("hab", "--version").Output()
	if err == nil {
		r := regexp.MustCompile(`hab (\d+\.\d+\.\d+)/`)
		currentVersion = r.FindStringSubmatch(string(output))[1]
	}

	if currentVersion == rootCmdOpts.version {
		cmd.Printf("Chef Habitat is already up-to-date (%s)\n", rootCmdOpts.version)
	} else {
		cmd.Printf("Going to install the %s build of Chef Habitat %s\n", viper.Get("target"), rootCmdOpts.version)

		installFileName := filepath.Join(ciutils.SettingsPath(HabitatInstallScript))

		installFileExists, err := fs.Exists(installFileName)
		if err != nil {
			return errors.Wrapf(err, "failed to determine if %s exists", installFileName)
		}

		if !installFileExists {
			if err = fs.DownloadRemoteFile(HabitatInstallScriptURL, installFileName); err != nil {
				return releaseLock(lock, err, "failed to download Chef Habitat install file")
			}

			if err = fs.Chmod(installFileName, 0777); err != nil {
				return releaseLock(lock, err, "failed to make temporary install file executable")
			}
		}

		if err = doInstallHabitat(installFileName); err != nil {
			return releaseLock(lock, err, "failed to install Chef Habitat")
		}
	}

	return releaseLock(lock, nil, "")
}

func releaseLock(lock filelock.TryLockerSafe, upstreamErr error, upstreamReason string) error {
	wrappedUpstreamErr := errors.Wrap(upstreamErr, upstreamReason)

	if err := lock.Unlock(); err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrap(err, "failed to release install lock"))
	}

	return wrappedUpstreamErr
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set environment variable prefix, eg: INSTALL_HABITAT_TARGET
	viper.SetEnvPrefix("install_habitat")

	// Load the config file
	settingsFile := ciutils.SettingsPath("install-habitat.toml")
	settingsFileExists, err := fs.Exists(settingsFile)
	if err != nil {
		log.Fatal(err)
	}

	if settingsFileExists {
		viper.SetConfigFile(settingsFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
	}

	// Override our config with any matching environment variables
	viper.AutomaticEnv()
}
