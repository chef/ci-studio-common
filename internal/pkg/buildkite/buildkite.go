package buildkite

import (
	"fmt"
	"path/filepath"

	"github.com/chef/ci-studio-common/internal/pkg/paths"
)

// AgentHook returns the path to the given hook on disk
func AgentHook(hookName string) string {
	return filepath.Join(AgentHooksDir, hookName)
}

// CIStudioCommonAgentHook returns the path to the given hook on disk
func CIStudioCommonAgentHook(hookName string) string {
	return filepath.Join(paths.InstallDir, "buildkite-agent-hooks", fmt.Sprintf("%s.sh", hookName))
}
