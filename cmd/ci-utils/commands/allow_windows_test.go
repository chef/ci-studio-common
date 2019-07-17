// +build windows

package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAllowWindows(t *testing.T) {
	var err error

	cmd := &cobra.Command{}
	args := []string{"buildkite-agent"}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	err = allowE(cmd, args)
	assert.NotNil(t, err)
	assert.Equal(t, "this command is not supported on Windows", err.Error())
}
