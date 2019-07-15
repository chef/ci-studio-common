package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"

	"github.com/chef/ci-studio-common/lib"

	"github.com/juju/fslock"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "install-habitat",
		Short: "Install VERSION of Chef Habitat from CHANNEL.",
		Run:   maybeInstallHabitat,
	}

	rootOpts = struct {
		version string
		channel string
		target  string
	}{}

	globalHabValue = "0.82.0"
)

func init() {
	defaultVersion := lib.SettingWithDefault("hab-version", globalHabValue)
	defaultTarget := lib.SettingWithDefault("hab-target", fmt.Sprintf("x86_64-%s", runtime.GOOS))

	rootCmd.Flags().StringVarP(&rootOpts.version, "version", "v", defaultVersion, "Which version of Habitat you wish to install.")
	rootCmd.Flags().StringVarP(&rootOpts.channel, "channel", "c", "stable", "The channel from which you wish to install Habitat.")
	rootCmd.Flags().StringVarP(&rootOpts.target, "target", "t", defaultTarget, "The kernel target for this installation.")
}

// Execute runs the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func maybeInstallHabitat(cmd *cobra.Command, args []string) {
	lock := fslock.New(lib.LockPath("install-habitat"))
	lockErr := lock.TryLock()

	if lockErr == nil {
		r := regexp.MustCompile(`hab (\d+\.\d+\.\d+)/`)

		currentVersion := "none"
		output, err := lib.ShellOut("hab", "--version").Output()
		if err == nil {
			currentVersion = r.FindStringSubmatch(string(output))[1]
		}

		if currentVersion == rootOpts.version {
			fmt.Println("Chef Habitat is already up-to-date")
		} else {
			doInstallHabitat()
		}

		err = lock.Unlock()
		lib.Check(err)
	} else {
		fmt.Println("Chef Habitat install already in progress -- skipping")
	}

}

func doInstallHabitat() {
	fmt.Printf("Going to install the %s build of Chef Habitat %s\n", rootOpts.target, rootOpts.version)

	writeOutUserFiles()

	if runtime.GOOS == "windows" {
		installHabitatChocolatey()
	} else {
		installHabitatScript()
	}
}

func writeOutUserFiles() {
	err := ioutil.WriteFile(lib.SettingsPath("hab-version"), []byte(rootOpts.version), 0644)
	lib.Check(err)

	err = ioutil.WriteFile(lib.SettingsPath("hab-target"), []byte(rootOpts.target), 0644)
	lib.Check(err)
}

func installHabitatScript() {
	installFile, err := ioutil.TempFile("", "hab-install")
	lib.Check(err)
	defer os.Remove(installFile.Name())

	downloadURL := "https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.sh"
	err = lib.DownloadFile(installFile.Name(), downloadURL)
	lib.Check(err)

	err = os.Chmod(installFile.Name(), 0777)
	lib.Check(err)
	installFile.Close()

	installCmd := lib.ShellOut(installFile.Name(), "-v", rootOpts.version, "-c", rootOpts.channel, "-t", rootOpts.target)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	err = installCmd.Run()
	lib.Check(err)
}

func installHabitatChocolatey() {
	installCmd := lib.ShellOut("choco", "install", "habitat", "--allow-downgrade", "-y", "--version", rootOpts.version)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	err := installCmd.Run()
	lib.Check(err)
}
