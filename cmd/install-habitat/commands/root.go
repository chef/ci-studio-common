package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"runtime"

	"github.com/juju/fslock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/chef/ci-studio-common/internal/pkg/paths"
	"github.com/chef/ci-studio-common/internal/pkg/system"
)

var (
	rootCmd = &cobra.Command{
		Use:   "install-habitat",
		Short: "Install VERSION of Chef Habitat from CHANNEL.",
		RunE:  maybeInstallHabitat,
	}

	rootOpts = struct {
		version string
		channel string
		target  string
	}{}

	// DO NOT MODIFY -- This value is automatically updated by Expeditor
	globalHabValue = "0.82.0"
)

func init() {
	defaultVersion := paths.SettingWithDefault("hab-version", globalHabValue)
	defaultTarget := paths.SettingWithDefault("hab-target", fmt.Sprintf("x86_64-%s", runtime.GOOS))

	rootCmd.Flags().StringVarP(&rootOpts.version, "version", "v", defaultVersion, "Which version of Habitat you wish to install.")
	rootCmd.Flags().StringVarP(&rootOpts.channel, "channel", "c", "stable", "The channel from which you wish to install Habitat.")
	rootCmd.Flags().StringVarP(&rootOpts.target, "target", "t", defaultTarget, "The kernel target for this installation.")
}

// Execute handles the execution of child commands and flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func maybeInstallHabitat(cmd *cobra.Command, args []string) error {
	lock := fslock.New(paths.LockPath("install-habitat"))
	lockErr := lock.TryLock()

	if lockErr == nil {
		r := regexp.MustCompile(`hab (\d+\.\d+\.\d+)/`)

		currentVersion := "none"
		output, err := system.ShellOut("hab", "--version").Output()
		if err == nil {
			currentVersion = r.FindStringSubmatch(string(output))[1]
		}

		if currentVersion == rootOpts.version {
			fmt.Println("Chef Habitat is already up-to-date")
		} else {
			fmt.Printf("Going to install the %s build of Chef Habitat %s\n", rootOpts.target, rootOpts.version)

			err := ioutil.WriteFile(paths.SettingsPath("hab-version"), []byte(rootOpts.version), 0644)
			if err != nil {
				return errors.Wrap(err, "failed to write out hab-version file")
			}

			err = ioutil.WriteFile(paths.SettingsPath("hab-target"), []byte(rootOpts.target), 0644)
			if err != nil {
				return errors.Wrap(err, "failed to write out hab-target file")
			}

			err = doInstallHabitat()
			if err != nil {
				return releaseLock(lock, err, "failed to install Chef Habitat")
			}
		}

		return releaseLock(lock, nil, "")
	} else {
		fmt.Println("Chef Habitat install already in progress -- skipping")
	}

	return nil
}

func releaseLock(lock *fslock.Lock, upstreamErr error, upstreamReason string) error {
	wrappedUpstreamErr := errors.Wrap(upstreamErr, upstreamReason)

	err := lock.Unlock()
	if err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrap(err, "failed to release install lock"))
	}

	return wrappedUpstreamErr
}
