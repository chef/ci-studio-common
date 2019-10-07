// +build !windows

package commands

import (
	"github.com/spf13/viper"
)

const (
	HabitatInstallScript    string = "install-habitat.sh"
	HabitatInstallScriptURL string = "https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.sh"
)

func doInstallHabitat(installFileName string) error {
	return execCommand(installFileName, "-v", rootCmdOpts.version, "-c", rootCmdOpts.channel, "-t", viper.GetString("target")).Run()
}
