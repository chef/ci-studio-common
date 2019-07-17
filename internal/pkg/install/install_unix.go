// +build !windows

package install

const DefaultRootParentDir string = "/opt"

const DefaultSettingsDir string = "/var/opt/ci-utils"

// binName returns the name of binary based on the OS
func binName(binary string) string {
	return binary
}
