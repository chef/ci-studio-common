package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/buildkite"
	"github.com/chef/ci-studio-common/internal/pkg/datastructs"
	"github.com/chef/ci-studio-common/internal/pkg/files"
)

var (
	hookCmd = &cobra.Command{
		Use:   "hook TYPE HOOK",
		Short: "Install one of the supported HOOKS as an Buildkite Agent Hook.",
		Long: `SUPPORTED HOOK TYPES:
  https://buildkite.com/docs/agent/v3/hooks#available-hooks

SUPPORTED HOOKS:
  ci-studio-common                Helper to run as an environment hook to load the ci-studio-common helpers.
  short-checkout-path             Checkout code into shorter code to deal with file path limits.`,
		Args: validHookArgs,
		RunE: installHook,
	}

	validHookTypes = []string{
		"environment",
		"pre-checkout", "checkout", "post-checkout",
		"pre-command", "command", "post-command",
		"pre-artifact", "post-artifact",
		"pre-exit",
	}

	validHookNames = []string{
		"ci-studio-common",
		"short-checkout-path",
	}
)

func init() {
	rootCmd.AddCommand(hookCmd)
}

func validHookArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("you must provide a hook type and a hook name")
	}

	if !datastructs.StringInSlice(args[0], validHookTypes) {
		return errors.Errorf("%s is not a valid Buildkite hook type", args[0])
	}

	if !datastructs.StringInSlice(args[1], validHookNames) {
		return errors.Errorf("%s is not a supported hook", args[1])
	}

	return nil
}

func installHook(cmd *cobra.Command, args []string) error {
	hookType := args[0]
	hookName := args[1]

	err := os.MkdirAll(buildkite.AgentHooksDir, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create agent hooks directory")
	}

	hookFile := buildkite.CIStudioCommonAgentHook(hookName)
	if files.FileExists(hookFile) {
		err = installShellHook(hookType, hookName)
		if err != nil {
			return errors.Wrapf(err, "failed to install the %s %s hook", hookName, hookType)
		}
	} else {
		return errors.Errorf("could not find hook %s", hookName)
	}

	return nil
}

func installShellHook(hookType string, hookName string) error {
	sourceFile := buildkite.CIStudioCommonAgentHook(hookName)
	footer := fmt.Sprintf("echo \"hook complete: %s\"", hookName)

	hookFilePath := buildkite.AgentHook(hookType)
	hookFile, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open the %s hook", hookFilePath)
	}
	defer hookFile.Close()

	hookFileContents, err := ioutil.ReadFile(hookFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read the contents of %s", hookFilePath)
	}

	fi, err := hookFile.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to get file stats for %s", hookFilePath)
	}

	if fi.Size() == 0 {
		_, err := hookFile.WriteString("#!/bin/bash\n\nset -eou pipefail\n\n")
		if err != nil {
			return errors.Wrapf(err, "failed to initialize contents of %s", hookFilePath)
		}
	}

	if !strings.Contains(string(hookFileContents), footer) {
		header := fmt.Sprintf("echo \"--- executing hook: %s\"", hookName)
		content := fmt.Sprintf(". %q", sourceFile)

		fullContents := fmt.Sprintf("%s\n%s\n%s\n", header, content, footer)
		_, err := hookFile.WriteString(fullContents)
		if err != nil {
			return errors.Wrapf(err, "failed to insert %s hook into %s", hookName, hookFilePath)
		}

		err = os.Chmod(hookFile.Name(), 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to update %s file permissions", hookFilePath)
		}
	}

	return nil
}
