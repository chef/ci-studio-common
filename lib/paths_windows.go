// +build windows

package lib

import (
	"fmt"
)

// InstallDirParent returns the parent directory of the installation (for tar unarchive)
const InstallDirParent string = "C:\\"

// SettingsDir returns the path to directory where settings are kept
const SettingsDir string = "C:\\ci-studio-settings"

// BinName returns the name of binary based on the OS
func BinName(binary string) string {
	return fmt.Sprintf("%s.exe", binary)
}
