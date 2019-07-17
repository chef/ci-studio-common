// +build !windows

package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func uninstallHabitat(cmd *cobra.Command, args []string) error {
	habBinPath, err := exec.LookPath("hab")

	if err == nil {
		err = os.Remove(habBinPath)
		if err != nil {
			return errors.Wrapf(err, "failed to remove %s during Chef Habitat uninstall", habBinPath)
		}

		err = os.RemoveAll("/hab")
		if err != nil {
			return errors.Wrap(err, "failed to remove /hab during Chef Habitat uninstall")
		}

		fmt.Println("Chef Habitat has been removed.")
	} else {
		fmt.Println("Chef Habitat has already been uninstalled.")
	}

	return nil
}
