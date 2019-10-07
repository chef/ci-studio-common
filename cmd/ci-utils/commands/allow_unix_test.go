// +build !windows

package commands

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
)

func TestAllowUnix(t *testing.T) {
	var err error

	fs = filesystem.NewMemFs()

	cmd := &cobra.Command{}
	args := []string{"buildkite-agent"}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	err = allowE(cmd, args)
	assert.Nil(t, err)

	expected := heredoc.Doc(`
		buildkite-agent ALL=NOPASSWD:SETENV: /bin/hab
		buildkite-agent ALL=NOPASSWD:SETENV: /usr/bin/ci-utils
		buildkite-agent ALL=NOPASSWD:SETENV: /usr/bin/install-habitat
	`)

	actual, err := fs.ReadFile("/etc/sudoers.d/buildkite-agent")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, string(actual))
}
