package lib

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

// InstallDirParent returns the parent directory of the installation (for tar unarchive)
func InstallDirParent() string {
	if runtime.GOOS == "windows" {
		return "C:\\"
	}

	return "/opt"
}

// InstallDir returns the path to the installation of ci-studio-common
func InstallDir() string {
	return filepath.Join(InstallDirParent(), "ci-studio-common")
}

// InstallBinPath returns the path to the binary
func InstallBinPath(binary string) string {
	return filepath.Join(InstallDir(), "bin", BinName(binary))
}

// InstallBackupDir returns the path where the current installation is backed up during upgrade
func InstallBackupDir() string {
	return filepath.Join(InstallDirParent(), "ci-studio-common-bak")
}

// InstallBackupBinPath returns the path to the binary in the backup directory
func InstallBackupBinPath(binary string) string {
	return filepath.Join(InstallBackupDir(), "bin", BinName(binary))
}

// BinName returns the name of binary based on the OS
func BinName(binary string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s.exe", binary)
	}

	return binary
}

// SettingsDir returns the path to directory where settings are kept
func SettingsDir() string {
	if runtime.GOOS == "windows" {
		return "C:\\ci-studio-settings"
	}

	return "/var/opt/ci-studio-common"
}

// SettingsPath will return the path to the settings file based on the installation
func SettingsPath(file string) string {
	return filepath.Join(SettingsDir(), file)
}

// LockPath will return the path to the lock file based on the installation
func LockPath(file string) string {
	return filepath.Join(SettingsDir(), fmt.Sprintf("%s.lock", file))
}

// SettingWithDefault will fetch the setting from file, otherwise return the default value
func SettingWithDefault(file string, defValue string) string {
	settingValue, err := ioutil.ReadFile(SettingsPath(file))
	if err != nil {
		settingValue = []byte(defValue)
	}

	return strings.TrimSpace(string(settingValue))
}
