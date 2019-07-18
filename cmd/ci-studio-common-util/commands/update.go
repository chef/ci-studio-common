package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/fslock"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/chef/ci-studio-common/internal/pkg/files"
	"github.com/chef/ci-studio-common/internal/pkg/http"
	"github.com/chef/ci-studio-common/internal/pkg/paths"
	"github.com/chef/ci-studio-common/internal/pkg/system"
)

var (
	updateCommand = &cobra.Command{
		Use:   "update",
		Short: "Update the ci-studio-common install",
		RunE:  maybeUpdateInstall,
	}

	updateOpts = struct {
		phase  int
		force  bool
		suffix string
	}{}

	assetOnDisk = filepath.Join(paths.InstallDirParent, "ci-studio-common.tar.gz")

	upgradeStatusFile = paths.SettingsPath("upgrade-in-progress")
)

func init() {
	rootCmd.AddCommand(updateCommand)

	updateCommand.Flags().IntVar(&updateOpts.phase, "phase", 1, "The phase of the update to complete")
	updateCommand.Flags().BoolVar(&updateOpts.force, "force", false, "Perform the installation even if no updates are available")
	updateCommand.Flags().StringVar(&updateOpts.suffix, "suffix", "rc", "The suffix to use when downloading the asset")
}

// When we run these phases, each phase triggers the next using cmd.Start(). We use Start() because it does
// not wait for the command to finish, which is what we want for this sequence to work. We do it this way
// because of limitations on Windows for binaries modifying themselves.
func maybeUpdateInstall(cmd *cobra.Command, args []string) error {
	switch updateOpts.phase {
	case 1:
		lock := fslock.New(paths.LockPath("upgrade-ci-studio-common"))
		lockErr := lock.TryLock()

		if lockErr == nil {
			return performPhaseOne()
		} else {
			fmt.Println("ci-studio-common upgrade already in progress -- skipping")
		}
	case 2:
		return performPhaseTwo()
	case 3:
		return performPhaseThree()
	default:
		return errors.Errorf("%d is an unsupported phase", updateOpts.phase)
	}

	return nil
}

// In phase one of the upgrade, we:
//   1. check if an upgrade is required
//   2. download the asset
//   3. copy the existing directory into a backup directory
//   4. begin phase two from the backup directory
func performPhaseOne() error {
	assetURL := fmt.Sprintf("https://chef-cd-artifacts.s3-us-west-2.amazonaws.com/ci-studio-common/ci-studio-common-2.0.0-%s-%s.tar.gz", runtime.GOOS, updateOpts.suffix)

	localEtag := paths.SettingWithDefault("etag", "none")
	remoteHeaders, err := http.GetURLHeaders(assetURL)
	if err != nil {
		return errors.Wrap(err, "failed to fetch etag header for new asset")
	}

	remoteEtag := remoteHeaders.Get("ETag")

	if (remoteEtag != localEtag) || updateOpts.force {
		phasePrintln(1, "Upgrade available")
		err := os.Chdir(paths.InstallDirParent)
		if err != nil {
			return errors.Wrap(err, "failed to cd into install directory parent")
		}

		// Mark that an upgrade is in progress
		err = ioutil.WriteFile(upgradeStatusFile, []byte(""), 0644)
		if err != nil {
			return errors.Wrap(err, "failed to create upgrade status file")
		}

		// Save new etag to disk
		err = ioutil.WriteFile(paths.SettingsPath("etag-new"), []byte(remoteEtag), 0644)
		if err != nil {
			return releaseUpgradeState(err, "failed to write new etag to disk")
		}

		// Download the asset to disk
		phasePrintln(1, "Downloading latest release")
		err = files.DownloadFile(assetOnDisk, assetURL)
		if err != nil {
			return releaseUpgradeState(err, "failed to download asset tarball")
		}

		// Try and backup the current installation
		phasePrintln(1, "Release downloaded -- backuping up current installation")
		err = files.CopyDir(paths.InstallDir, paths.InstallBackupDir, true)
		if err != nil {
			return releaseUpgradeState(err, "failed to backup existing installation")
		}

		// Begin phase two
		cmd := system.ShellOut(paths.InstallBackupBinPath("ci-studio-common-util"), "update", "--phase", "2")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			return releaseUpgradeState(err, "failed to launch 2nd phase of upgrade")
		}
	} else {
		fmt.Println("ci-studio-common is up-to-date")
	}

	return nil
}

// In phase two of the upgrade, we:
//   1. remove the old install directory
//   2. move the new install into place
//   3. begin phase three from the new install
func performPhaseTwo() error {
	err := os.Chdir(paths.InstallDirParent)
	if err != nil {
		return releaseUpgradeState(err, "failed to cd into install directory parent")
	}

	// Remove the previous installation
	phasePrintln(2, "Removing current installation")
	err = os.RemoveAll(paths.InstallDir)
	if err != nil {
		return releaseUpgradeState(err, "failed to remove existing installation")
	}

	// Untar the new install into place
	phasePrintln(2, "Moving new installation into place")
	err = archiver.Unarchive(assetOnDisk, paths.InstallDirParent)
	if err != nil {
		return releaseUpgradeState(err, "failed to untar new installation into place")
	}

	// Begin phase three
	cmd := system.ShellOut(paths.InstallBinPath("ci-studio-common-util"), "update", "--phase", "3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return releaseUpgradeState(err, "failed to launch 3rd phase of upgrade")
	}

	return nil
}

// In phase three of the upgrade, we:
//   1. remove the backup directory
func performPhaseThree() error {
	err := os.Chdir(paths.InstallDirParent)
	if err != nil {
		return releaseUpgradeState(err, "failed to cd into install directory parent")
	}

	phasePrintln(3, "Upgrade complete -- cleaning up")

	// Move the new etag into place
	err = os.RemoveAll(paths.SettingsPath("etag"))
	if err != nil {
		return releaseUpgradeState(err, "failed to remove existing etag from disk")
	}

	err = files.RenameFile(paths.SettingsPath("etag-new"), paths.SettingsPath("etag"))
	if err != nil {
		return releaseUpgradeState(err, "failed to move new etag into place on disk")
	}

	return releaseUpgradeState(nil, "")
}

// phasePrintln provides a consistent mechanism for printing upgrade state through the different phases
func phasePrintln(phase int, message string) {
	fmt.Printf("--> {%d} %s", phase, message)
}

func releaseUpgradeState(upstreamErr error, upstreamReason string) error {
	wrappedUpstreamErr := errors.Wrap(upstreamErr, upstreamReason)

	err := os.RemoveAll(paths.InstallBackupDir)
	if err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s", paths.InstallBackupDir))
	}

	err = os.RemoveAll(assetOnDisk)
	if err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s from disk", assetOnDisk))
	}

	err = os.RemoveAll(upgradeStatusFile)
	if err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s from disk", upgradeStatusFile))
	}

	return errors.Wrap(upstreamErr, upstreamReason)
}
