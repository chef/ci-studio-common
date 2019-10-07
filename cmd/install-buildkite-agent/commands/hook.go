package commands

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/buildkite"
)

var (
	hookCmd = &cobra.Command{
		Use:   "hook TYPE HOOK",
		Short: "Install one of the supported HOOKS as an Buildkite Agent Hook.",
		Long: `SUPPORTED HOOK TYPES:
  https://buildkite.com/docs/agent/v3/hooks#available-hooks

SUPPORTED HOOKS:
  ci-utils                Helper to run as an environment hook to load the ci-utils helpers.
  short-checkout-path     Checkout code into shorter code to deal with file path limits.`,
		Args: validHookArgs,
		RunE: hookE,
	}

	validHookTypes = []string{
		"environment",
		"pre-checkout", "checkout", "post-checkout",
		"pre-command", "command", "post-command",
		"pre-artifact", "post-artifact",
		"pre-exit",
	}

	validHookNames = []string{
		"ci-utils",
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

	if !stringInSlice(args[0], validHookTypes) {
		return errors.Errorf("%s is not a valid Buildkite hook type", args[0])
	}

	if !stringInSlice(args[1], validHookNames) {
		return errors.Errorf("%s is not a supported hook", args[1])
	}

	return nil
}

func hookE(cmd *cobra.Command, args []string) error {
	hookType := args[0]
	hookName := args[1]

	err := fs.MkdirAll(buildkite.AgentHooksDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", buildkite.AgentHooksDir)
	}

	hookFile := ciutils.AgentHook(hookName)
	hookFileExists, err := fs.Exists(hookFile)
	if err != nil {
		return errors.Wrapf(err, "failed to determine if file %s exists", hookFile)
	}

	if hookFileExists {
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
	// Make sure the buildkite hook file exists
	bkHookFilePath := buildkite.AgentHook(hookType)
	bkHookFileExists, err := fs.Exists(bkHookFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to determine if %s exists", bkHookFilePath)
	}

	if !bkHookFileExists {
		err := fs.WriteFile(bkHookFilePath, []byte("#!/bin/bash\n\nset -eou pipefail\n"), 0755)
		if err != nil {
			return errors.Wrapf(err, "failed to initialize %s", bkHookFilePath)
		}
	}

	// Generate the content we're going to inject and inject it if necessary
	hookContents := heredoc.Docf(`

		echo "--- executing hook: %s"
		. %s
		echo "hook complete: %s"`, hookName, ciutils.AgentHook(hookName), hookName)
	hookContentsBytes := []byte(hookContents)

	hookPresent, err := fs.FileContainsBytes(bkHookFilePath, hookContentsBytes)
	if err != nil {
		return errors.Wrapf(err, "failed to determine if hook %s is already present in %s", hookName, bkHookFilePath)
	}

	if hookPresent {
		return nil
	}

	return fs.AppendIfMissing(bkHookFilePath, hookContentsBytes, 0755)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
