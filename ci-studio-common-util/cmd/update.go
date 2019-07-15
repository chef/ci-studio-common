package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chef/ci-studio-common/lib"
	"github.com/juju/fslock"
	"github.com/mholt/archiver"

	"github.com/spf13/cobra"
)

var (
	updateCommand = &cobra.Command{
		Use:   "update",
		Short: "Update the ci-studio-common install",
		Run:   maybeUpdateInstall,
	}

	updateOpts = struct {
		phase  int
		force  bool
		suffix string
	}{}

	assetOnDisk = filepath.Join(lib.InstallDirParent, "ci-studio-common.tar.gz")

	upgradeStatusFile = lib.SettingsPath("upgrade-in-progress")
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
func maybeUpdateInstall(cmd *cobra.Command, args []string) {
	switch updateOpts.phase {
	case 1:
		lock := fslock.New(lib.LockPath("upgrade-ci-studio-common"))
		lockErr := lock.TryLock()

		if lockErr == nil {
			performPhaseOne()
		} else {
			fmt.Println("ci-studio-common upgrade already in progress -- skipping")
		}
	case 2:
		performPhaseTwo()
	case 3:
		performPhaseThree()
	default:
		log.Fatal("Unsupported phase")
	}
}

// In phase one of the upgrade, we:
//   1. check if an upgrade is required
//   2. download the asset
//   3. copy the existing directory into a backup directory
//   4. begin phase two from the backup directory
func performPhaseOne() {
	assetURL := fmt.Sprintf("https://chef-cd-artifacts.s3-us-west-2.amazonaws.com/ci-studio-common/ci-studio-common-2.0.0-%s-%s.tar.gz", runtime.GOOS, updateOpts.suffix)

	localEtag := lib.SettingWithDefault("etag", "none")
	remoteEtag := lib.GetURLHeaderByKey(assetURL, "ETag")

	if (remoteEtag != localEtag) || updateOpts.force {
		fmt.Println("--> {1} Upgrade available")
		err := os.Chdir(lib.InstallDirParent)
		lib.Check(err)

		// Mark that an upgrade is in progress
		err = ioutil.WriteFile(upgradeStatusFile, []byte(""), 0644)
		lib.Check(err)

		// Save new etag to disk
		err = ioutil.WriteFile(lib.SettingsPath("etag-new"), []byte(remoteEtag), 0644)
		lib.Check(err)

		// Download the asset to disk
		fmt.Println("--> {1} Downloading latest release")
		err = lib.DownloadFile(assetOnDisk, assetURL)
		lib.Check(err)

		// Try and backup the current installation
		fmt.Println("--> {1} Release downloaded -- backuping up current installation")
		err = lib.CopyDir(lib.InstallDir(), lib.InstallBackupDir(), true)
		lib.Check(err)

		// Begin phase two
		cmd := lib.ShellOut(lib.InstallBackupBinPath("ci-studio-common-util"), "update", "--phase", "2")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Start()
		lib.Check(err)
	} else {
		fmt.Println("ci-studio-common is up-to-date")
	}
}

// In phase two of the upgrade, we:
//   1. remove the old install directory
//   2. move the new install into place
//   3. begin phase three from the new install
func performPhaseTwo() {
	err := os.Chdir(lib.InstallDirParent)
	lib.Check(err)

	// Remove the previous installation
	fmt.Println("--> {2} Removing current installation")
	err = os.RemoveAll(lib.InstallDir())
	lib.Check(err)

	// Untar the new install into place
	fmt.Println("--> {2} Moving new installation into place")
	err = archiver.Unarchive(assetOnDisk, lib.InstallDirParent)
	lib.Check(err)

	// Begin phase three
	cmd := lib.ShellOut(lib.InstallBinPath("ci-studio-common-util"), "update", "--phase", "3")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	lib.Check(err)
}

// In phase three of the upgrade, we:
//   1. remove the backup directory
func performPhaseThree() {
	err := os.Chdir(lib.InstallDirParent)
	lib.Check(err)

	fmt.Println("--> {3} Upgrade complete -- cleaning up")

	// Move the new etag into place
	err = os.RemoveAll(lib.SettingsPath("etag"))
	lib.Check(err)

	err = lib.RenameFile(lib.SettingsPath("etag-new"), lib.SettingsPath("etag"))
	lib.Check(err)

	// Cleanup the backup directory
	err = os.RemoveAll(lib.InstallBackupDir())
	lib.Check(err)

	// Cleanup the tarball
	err = os.RemoveAll(assetOnDisk)
	lib.Check(err)

	err = os.RemoveAll(upgradeStatusFile)
	lib.Check(err)
}
