// +build !windows

package paths

// InstallDirParent returns the parent directory of the installation (for tar unarchive)
const InstallDirParent string = "/opt"

// InstallDir returns the path to the installation of ci-studio-common
const InstallDir string = "/opt/ci-studio-common"

// InstallBackupDir returns the path where the current installation is backed up during upgrade
const InstallBackupDir string = "/opt/ci-studio-common-bak"

// SettingsDir returns the path to directory where settings are kept
const SettingsDir string = "/var/opt/ci-studio-common"

// BinName returns the name of binary based on the OS
func BinName(binary string) string {
	return binary
}
