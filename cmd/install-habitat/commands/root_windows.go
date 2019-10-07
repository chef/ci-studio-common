// +build windows

package commands

import (
	"fmt"
)

const (
	HabitatInstallScript    string = "install-habitat.ps1"
	HabitatInstallScriptURL string = "https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.ps1"
)

func doInstallHabitat(installFileName string) error {
	return execCommand("powershell.exe", "-Command", fmt.Sprintf("%s -v %s -c %s", installFileName, rootCmdOpts.version, rootCmdOpts.channel)).Run()
}
