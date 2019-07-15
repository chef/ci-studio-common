package paths

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// InstallBinPath returns the path to the binary
func InstallBinPath(binary string) string {
	return filepath.Join(InstallDir, "bin", BinName(binary))
}

// InstallBackupBinPath returns the path to the binary in the backup directory
func InstallBackupBinPath(binary string) string {
	return filepath.Join(InstallBackupDir, "bin", BinName(binary))
}

// SettingsPath will return the path to the settings file based on the installation
func SettingsPath(file string) string {
	return filepath.Join(SettingsDir, file)
}

// LockPath will return the path to the lock file based on the installation
func LockPath(file string) string {
	return filepath.Join(SettingsDir, fmt.Sprintf("%s.lock", file))
}

// SettingWithDefault will fetch the setting from file, otherwise return the default value
func SettingWithDefault(file string, defValue string) string {
	settingValue, err := ioutil.ReadFile(SettingsPath(file))
	if err != nil {
		settingValue = []byte(defValue)
	}

	return strings.TrimSpace(string(settingValue))
}
