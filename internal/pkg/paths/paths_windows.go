// +build windows

package paths

import (
	"fmt"
)

// InstallDirParent returns the parent directory of the installation (for tar unarchive)
const InstallDirParent string = "C:\\"

// InstallDir returns the path to the installation of ci-studio-common
const InstallDir string = "C:\\ci-studio-common"

// InstallBackupDir returns the path where the current installation is backed up during upgrade
const InstallBackupDir string = "C:\\ci-studio-common-bak"

// SettingsDir returns the path to directory where settings are kept
const SettingsDir string = "C:\\ci-studio-settings"

// BinName returns the name of binary based on the OS
func BinName(binary string) string {
	return fmt.Sprintf("%s.exe", binary)
}
