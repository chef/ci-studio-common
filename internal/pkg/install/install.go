package install

import (
	"fmt"
	"path/filepath"
)

// Install - interface for install.
type Install interface {
	AgentHook(hookName string) string
	BackupDir() string
	BackupBinPath(binary string) string
	BinPath(binary string) string
	LockPath(file string) string
	RootDir() string
	RootDirParent() string
	SettingsPath(file string) string
}

// OsInstall - struct for OsInstall.
type OsInstall struct {
	RootParentDir string
	SettingsDir   string
}

// DefaultInstall returns the default installation interface.
func DefaultInstall() *OsInstall {
	return &OsInstall{
		RootParentDir: DefaultRootParentDir,
		SettingsDir:   DefaultSettingsDir,
	}
}

// AgentHook returns the path to the given hook on disk.
func (o *OsInstall) AgentHook(hookName string) string {
	return filepath.Join(o.RootDir(), "buildkite-agent-hooks", fmt.Sprintf("%s.sh", hookName))
}

// BackupDir returns the path where the current installation is backed up during upgrade.
func (o *OsInstall) BackupDir() string {
	return filepath.Join(o.RootParentDir, "ci-utils-bak")
}

// BackupBinPath returns the path to the binary in the backup directory.
func (o *OsInstall) BackupBinPath(binary string) string {
	return filepath.Join(o.BackupDir(), "bin", binName(binary))
}

// BinPath returns the path to the binary.
func (o *OsInstall) BinPath(binary string) string {
	return filepath.Join(o.RootDir(), "bin", binName(binary))
}

// LockPath will return the path to the lock file based on the installation.
func (o *OsInstall) LockPath(file string) string {
	return filepath.Join(o.SettingsDir, fmt.Sprintf("%s.lock", file))
}

// RootDir returns the path to the installation of ci-utils.
func (o *OsInstall) RootDir() string {
	return filepath.Join(o.RootParentDir, "ci-utils")
}

// RootDirParent returns the path to the root of the installation.
func (o *OsInstall) RootDirParent() string {
	return o.RootParentDir
}

// SettingsPath will return the path to the settings file based on the installation.
func (o *OsInstall) SettingsPath(file string) string {
	return filepath.Join(o.SettingsDir, file)
}
