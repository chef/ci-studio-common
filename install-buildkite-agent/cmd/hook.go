package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/chef/ci-studio-common/lib"

	"github.com/spf13/cobra"
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
		Args: cobra.MinimumNArgs(2),
		Run:  installHook,
	}
)

func init() {
	rootCmd.AddCommand(hookCmd)
}

func installHook(cmd *cobra.Command, args []string) {
	hookType := args[0]
	hookName := args[1]

	err := os.MkdirAll(lib.BuildkiteAgentHooksDir, 0755)
	lib.Check(err)

	hookFile := lib.CiStudioCommonAgentHook(hookName)
	if lib.FileExists(hookFile) {
		err = installShellHook(hookType, hookName)
		lib.Check(err)
	} else {
		log.Fatalf("Could not find the hook %s", hookName)
	}
}

func installShellHook(hookType string, hookName string) error {
	sourceFile := lib.CiStudioCommonAgentHook(hookName)
	footer := fmt.Sprintf("echo \"hook complete: %s\"", hookName)

	hookFilePath := lib.BuildkiteAgentHook(hookType)
	hookFile, err := os.OpenFile(hookFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	lib.Check(err)
	defer hookFile.Close()

	hookFileContents, err := ioutil.ReadFile(hookFilePath)
	lib.Check(err)

	fi, err := hookFile.Stat()
	lib.Check(err)

	if fi.Size() == 0 {
		_, err := hookFile.Write([]byte("#!/bin/bash\n\nset -eou pipefail\n\n"))
		lib.Check(err)
	}

	if !strings.Contains(string(hookFileContents), footer) {
		header := fmt.Sprintf("echo \"--- executing hook: %s\"", hookName)
		content := fmt.Sprintf(". %q", sourceFile)

		fullContents := fmt.Sprintf("%s\n%s\n%s\n", header, content, footer)
		_, err := hookFile.Write([]byte(fullContents))
		lib.Check(err)

		err = os.Chmod(hookFile.Name(), 0755)
		lib.Check(err)
	}

	return nil
}
