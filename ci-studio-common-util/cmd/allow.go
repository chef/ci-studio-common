package cmd

import (
	"fmt"
	"runtime"

	"github.com/chef/ci-studio-common/lib"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allowCmd)
}

var allowCmd = &cobra.Command{
	Use:   "allow USER",
	Short: "Allow USER to perform certain necessary operations with sudo.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			fmt.Println("This command is not currently supported on Windows.")
		} else {
			lib.AddSudoPermission("/bin/hab", args[0])
			lib.AddSudoPermission("/usr/bin/ci-studio-common-util", args[0])
			lib.AddSudoPermission("/usr/bin/install-habitat", args[0])
		}
	},
}
