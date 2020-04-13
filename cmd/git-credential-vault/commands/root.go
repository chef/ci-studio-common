package commands

import (
	"log"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

var (
	execCommand = exec.Command

	ciutils install.Install

	fs filesystem.FileSystem

	rootCmd = &cobra.Command{
		Use:              "git-credential-vault",
		Short:            "Fetch GitHub credentials from GitHub using the GitHub App secrets plugin.",
		SilenceUsage:     true,
		TraverseChildren: true,
	}

	secretsClient secrets.Client

	rootOpts = &rootOptions{}
)

type rootOptions struct {
	install string
}

// Execute handles the execution of child commands and flags
func Execute() {
	var err error

	fs = filesystem.NewOsFs()
	ciutils = install.DefaultInstall()

	secretsClient, err = secrets.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if err = rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("vault.dynamic_mount", "account/dynamic")
	viper.SetDefault("vault.static_mount", "account/static")

	rootCmd.Flags().StringVar(&rootOpts.install, "install", "", "GitHub App installation to use to generate credentials (required)")
	if err := rootCmd.MarkFlagRequired("install"); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Override our config with any matching environment variables
	viper.AutomaticEnv()

	// Set environment variable prefix, eg: VAULT_UTIL_AWS_REGION
	viper.SetEnvPrefix("vault_util")

	// Load the config file
	settingsFile := ciutils.SettingsPath("vault-util.toml")
	settingsFileExists, err := fs.Exists(settingsFile)
	if err != nil {
		log.Fatal(err)
	}

	if settingsFileExists {
		viper.SetConfigFile(settingsFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
	}

	// Override our config with any matching environment variables
	viper.AutomaticEnv()
}
