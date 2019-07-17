// +build windows

package install

import (
	"fmt"
)

const DefaultRootParentDir string = "C:\\"

const DefaultSettingsDir string = "C:\\ci-settings"

// binName returns the name of binary based on the OS
func binName(binary string) string {
	return fmt.Sprintf("%s.exe", binary)
}
