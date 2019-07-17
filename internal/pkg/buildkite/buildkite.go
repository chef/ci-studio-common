package buildkite

import (
	"path/filepath"
)

// AgentHook returns the path to the given hook on disk
func AgentHook(hookName string) string {
	return filepath.Join(AgentHooksDir, hookName)
}
