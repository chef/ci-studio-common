// +build !windows

package commands

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/chef/ci-studio-common/internal/pkg/files"
	"github.com/chef/ci-studio-common/internal/pkg/system"
)

func doInstallHabitat() error {
	installFile, err := ioutil.TempFile("", "hab-install")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary install.sh file")
	}
	defer os.Remove(installFile.Name())

	downloadURL := "https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.sh"
	err = files.DownloadFile(installFile.Name(), downloadURL)
	if err != nil {
		return errors.Wrap(err, "failed to download Chef Habitat install.sh file")
	}

	err = os.Chmod(installFile.Name(), 0777)
	if err != nil {
		return errors.Wrap(err, "failed to make temporary install.sh file executable")
	}
	installFile.Close()

	installCmd := system.ShellOut(installFile.Name(), "-v", rootOpts.version, "-c", rootOpts.channel, "-t", rootOpts.target)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	err = installCmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to install Chef Habitat")
	}

	return nil
}
