package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/chef/ci-studio-common/lib"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Completely uninstall Chef Habitat from the system.",
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			installCmd := lib.ShellOut("choco", "uninstall", "habitat")
			err := installCmd.Run()
			lib.Check(err)

			fmt.Println("Chef Habitat has been removed.")
		} else {
			habBinPath, err := exec.LookPath("hab")

			if err == nil {
				err = os.Remove(habBinPath)
				lib.Check(err)

				err = os.RemoveAll("/hab")
				lib.Check(err)

				fmt.Println("Chef Habitat has been removed.")
			} else {
				fmt.Println("Chef Habitat has already been uninstalled.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
