package commands

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
)

type updateCmdOptions struct {
	channel string
	force   bool
	phase   int
}

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the ci-utils install",
		RunE:  updateE,
	}

	updateCmdOpts = &updateCmdOptions{}
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&updateCmdOpts.channel, "channel", "stable", "Channel from which to download release")
	updateCmd.Flags().BoolVar(&updateCmdOpts.force, "force", false, "Perform the installation even if no updates are available")
	updateCmd.Flags().IntVar(&updateCmdOpts.phase, "phase", 1, "The phase of the update to complete")
}

// When we run these phases, each phase triggers the next using cmd.Start(). We use Start() because it does
// not wait for the command to finish, which is what we want for this sequence to work. We do it this way
// because of limitations on Windows for binaries modifying themselves.
func updateE(cmd *cobra.Command, args []string) error {
	switch updateCmdOpts.phase {
	case 1:
		lockfile := ciutils.LockPath("upgrade-ci-utils")
		lock, err := fslock.GetLock(lockfile)
		if err != nil {
			return errors.Wrapf(err, "could not create lockfile (%s)", lockfile)
		}

		_, err = lock.TryLock()
		if err == nil {
			return performPhaseOne(cmd)
		} else if err == filesystem.ErrLocked {
			cmd.Println("ci-utils upgrade already in progress -- skipping")
		} else {
			return errors.Wrap(err, "could not acquire file lock")
		}
	case 2:
		return performPhaseTwo(cmd)
	case 3:
		return performPhaseThree(cmd)
	default:
		return errors.Errorf("%d is an unsupported phase", updateCmdOpts.phase)
	}

	return nil
}

// In phase one of the upgrade, we:
//   1. check if an upgrade is required
//   2. download the asset
//   3. copy the existing directory into a backup directory
//   4. begin phase two from the backup directory
func performPhaseOne(cmd *cobra.Command) error {
	var localEtag []byte

	assetURL := fmt.Sprintf("https://packages.chef.io/files/%s/ci-utils/latest/ci-utils-%s.tar.gz", updateCmdOpts.channel, runtime.GOOS)
	localEtag, err := fs.ReadFile(ciutils.SettingsPath("etag"))
	if err != nil {
		localEtag = []byte("none")
	}

	// Fetch the remote etag
	response, err := http.Head(assetURL)
	if err != nil {
		return releaseUpgradeState(err, fmt.Sprintf("unable to download URL (%s)", assetURL))
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		statusErr := errors.Errorf("request unsuccessful: status code %s", response.Status)
		return releaseUpgradeState(statusErr, fmt.Sprintf("unable to download URL (%s)", assetURL))
	}

	remoteEtag := response.Header.Get("etag")

	// Determine if we need to perform an upgrade
	if (remoteEtag != string(localEtag)) || updateCmdOpts.force {
		phasePrintln(cmd, 1, "Upgrade available")

		if err := fs.WriteFile(ciutils.SettingsPath("upgrade-in-progress"), []byte(""), 0644); err != nil {
			return errors.Wrap(err, "failed to create upgrade status file")
		}

		if err := fs.WriteFile(ciutils.SettingsPath("etag-new"), []byte(remoteEtag), 0644); err != nil {
			return releaseUpgradeState(err, "failed to write new etag to disk")
		}

		phasePrintln(cmd, 1, "Downloading latest release")
		if err := fs.DownloadRemoteFile(assetURL, assetOnDisk()); err != nil {
			return releaseUpgradeState(err, "failed to download asset tarball")
		}

		phasePrintln(cmd, 1, "Release downloaded -- backing up current installation")
		if err := fs.CopyDir(ciutils.RootDir(), ciutils.BackupDir()); err != nil {
			return releaseUpgradeState(err, "failed to backup existing installation")
		}

		if err := execCommand(ciutils.BackupBinPath("ci-utils"), "update", "--phase", "2").Start(); err != nil {
			return releaseUpgradeState(err, "failed to launch 2nd phase of upgrade")
		}
	} else {
		cmd.Println("ci-utils is up-to-date")
	}

	return nil
}

// In phase two of the upgrade, we:
//   1. remove the old install directory
//   2. move the new install into place
//   3. begin phase three from the new install
func performPhaseTwo(cmd *cobra.Command) error {
	phasePrintln(cmd, 2, "Removing current installation")
	if err := fs.RemoveAll(ciutils.RootDir()); err != nil {
		return releaseUpgradeState(err, "failed to remove existing installation")
	}

	phasePrintln(cmd, 2, "Moving new installation into place")
	if err := archiver.Unarchive(assetOnDisk(), ciutils.RootDirParent()); err != nil {
		return releaseUpgradeState(err, "failed to untar new installation into place")
	}

	if err := execCommand(ciutils.BinPath("ci-utils"), "update", "--phase", "3").Start(); err != nil {
		return releaseUpgradeState(err, "failed to launch 3rd phase of upgrade")
	}

	return nil
}

// In phase three of the upgrade, we:
//   1. remove the backup directory
func performPhaseThree(cmd *cobra.Command) error {
	phasePrintln(cmd, 3, "Upgrade complete -- cleaning up")
	if err := fs.Rename(ciutils.SettingsPath("etag-new"), ciutils.SettingsPath("etag")); err != nil {
		return releaseUpgradeState(err, "failed to move new etag into place on disk")
	}

	return releaseUpgradeState(nil, "")
}

// phasePrintln provides a consistent mechanism for printing upgrade state through the different phases
func phasePrintln(cmd *cobra.Command, phase int, message string) {
	cmd.Printf("--> {%d} %s\n", phase, message)
}

// releaseUpgradeState handles cleaning up upgrade state files in the context of a possible error
func releaseUpgradeState(upstreamErr error, upstreamReason string) error {
	wrappedUpstreamErr := errors.Wrap(upstreamErr, upstreamReason)

	if err := fs.RemoveAll(ciutils.BackupDir()); err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s", ciutils.BackupDir()))
	}

	if err := fs.RemoveAll(assetOnDisk()); err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s from disk", assetOnDisk()))
	}

	if err := fs.RemoveAll(ciutils.SettingsPath("upgrade-in-progress")); err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrapf(err, "failed to remove %s from disk", ciutils.SettingsPath("upgrade-in-progress")))
	}

	lock, err := fslock.GetLock(ciutils.LockPath("upgrade-ci-utils"))
	if err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrap(err, "failed to find lock file"))
	}
	if err := lock.Unlock(); err != nil {
		return multierr.Append(wrappedUpstreamErr, errors.Wrap(err, "failed to release lock file"))
	}

	return errors.Wrap(upstreamErr, upstreamReason)
}

// assetOnDisk returns the path to where the asset exists on the local filesystem
func assetOnDisk() string {
	return filepath.Join(ciutils.RootDirParent(), "ci-utils.tar.gz")
}
