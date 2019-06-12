package lib

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// BuildkiteAgentHooksDir returns OS-specific path to the hooks directory
func BuildkiteAgentHooksDir() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "/usr/local/etc/buildkite-agent/hooks"
	case "windows":
		return "C:\\buildkite-agent\\hooks"
	default:
		return "/etc/buildkite-agent/hooks"
	}
}

// BuildkiteAgentHook returns the path to the given hook on disk
func BuildkiteAgentHook(hookName string) string {
	return filepath.Join(BuildkiteAgentHooksDir(), hookName)
}

// CiStudioCommonAgentHook returns the path to the given hook on disk
func CiStudioCommonAgentHook(hookName string) string {
	return filepath.Join(InstallDir(), "buildkite-agent-hooks", fmt.Sprintf("%s.sh", hookName))
}
