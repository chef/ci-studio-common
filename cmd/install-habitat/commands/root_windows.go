// +build windows

package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"

	"github.com/chef/ci-studio-common/internal/pkg/files"
	"github.com/chef/ci-studio-common/internal/pkg/system"
)

func doInstallHabitat() error {
	installFile, err := ioutil.TempFile("", "hab-install")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary install.ps1 file")
	}
	defer os.Remove(installFile.Name())

	downloadURL := "https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.ps1"
	err = files.DownloadFile(installFile.Name(), downloadURL)
	if err != nil {
		return errors.Wrap(err, "failed to download Chef Habitat install.ps1 file")
	}
	installFile.Close()

	powershellCommand := fmt.Sprintf("%s -v %s -c %s", installFile.Name(), rootOpts.version, rootOpts.channel)
	installCmd := system.ShellOut("powershell.exe", "-Command", powershellCommand)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	err = installCmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to install Chef Habitat")
	}

	return nil
}
