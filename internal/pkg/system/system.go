package system

import (
	"os/exec"
)

// ShellOut is a wrapper around exec.Command for easier mocking
func ShellOut(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
