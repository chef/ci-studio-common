package lib

import (
	"fmt"
	"path/filepath"
)

// BuildkiteAgentHook returns the path to the given hook on disk
func BuildkiteAgentHook(hookName string) string {
	return filepath.Join(BuildkiteAgentHooksDir, hookName)
}

// CiStudioCommonAgentHook returns the path to the given hook on disk
func CiStudioCommonAgentHook(hookName string) string {
	return filepath.Join(InstallDir(), "buildkite-agent-hooks", fmt.Sprintf("%s.sh", hookName))
}
