package commands

import (
	"bytes"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/chef/ci-studio-common/internal/pkg/buildkite"
	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
)

func TestValidHookArgs(t *testing.T) {
	cmd := &cobra.Command{}

	t.Run("not enough options", func(t *testing.T) {
		err := validHookArgs(cmd, []string{"ci-utils"})
		assert.NotNil(t, err)
		assert.Equal(t, "you must provide a hook type and a hook name", err.Error())
	})

	t.Run("invalid buildkite hook", func(t *testing.T) {
		err := validHookArgs(cmd, []string{"foobar", "ci-utils"})
		assert.NotNil(t, err)
		assert.Equal(t, "foobar is not a valid Buildkite hook type", err.Error())
	})

	t.Run("invalid ci-utils hook", func(t *testing.T) {
		err := validHookArgs(cmd, []string{"environment", "foobar"})
		assert.NotNil(t, err)
		assert.Equal(t, "foobar is not a supported hook", err.Error())
	})
}

func TestHookE(t *testing.T) {
	fs = filesystem.NewMemFs()
	ciutils = install.DefaultInstall()

	bkFileName := buildkite.AgentHook("environment")

	cmd := &cobra.Command{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	// stub out our hook contents
	err := fs.WriteFile(ciutils.AgentHook("ci-utils"), []byte("fake contents"), 0644)
	assert.Nil(t, err)

	t.Run("buildkite hook file does not exist", func(t *testing.T) {
		err = hookE(cmd, []string{"environment", "ci-utils"})
		assert.Nil(t, err)

		actual, err := fs.ReadFile(bkFileName)
		if err != nil {
			t.Fatal(err)
		}

		hookLocation := ciutils.AgentHook("ci-utils")
		expected := heredoc.Docf(`
			#!/bin/bash

			set -eou pipefail

			echo "--- executing hook: ci-utils"
			. %s
			echo "hook complete: ci-utils"
		`, hookLocation)

		assert.Equal(t, expected, string(actual))
	})

	t.Run("buildkite hook file does not contain hook", func(t *testing.T) {
		existing := heredoc.Doc(`
			#!/bin/bash

			set -eou pipefail

			. /some/internal/hook.sh
		`)

		hookLocation := ciutils.AgentHook("ci-utils")
		expected := heredoc.Docf(`
			#!/bin/bash

			set -eou pipefail

			. /some/internal/hook.sh

			echo "--- executing hook: ci-utils"
			. %s
			echo "hook complete: ci-utils"
		`, hookLocation)

		err := fs.WriteFile(bkFileName, []byte(existing), 0755)
		if err != nil {
			t.Fatal(err)
		}

		err = hookE(cmd, []string{"environment", "ci-utils"})
		assert.Nil(t, err)

		actual, err := fs.ReadFile(bkFileName)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, string(actual))
	})

	t.Run("buildkite hook file already contains hook", func(t *testing.T) {
		hookLocation := ciutils.AgentHook("ci-utils")
		existing := heredoc.Docf(`
			#!/bin/bash

			set -eou pipefail


			echo "--- executing hook: ci-utils"
			. %s
			echo "hook complete: ci-utils"
		`, hookLocation)

		err := fs.WriteFile(bkFileName, []byte(existing), 0755)
		if err != nil {
			t.Fatal(err)
		}

		err = hookE(cmd, []string{"environment", "ci-utils"})
		assert.Nil(t, err)

		actual, err := fs.ReadFile(bkFileName)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, existing, string(actual))
	})
}
