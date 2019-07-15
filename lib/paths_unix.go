// +build !windows

package lib

// InstallDirParent returns the parent directory of the installation (for tar unarchive)
const InstallDirParent string = "/opt"

// SettingsDir returns the path to directory where settings are kept
const SettingsDir string = "/var/opt/ci-studio-common"

// BinName returns the name of binary based on the OS
func BinName(binary string) string {
	return binary
}
