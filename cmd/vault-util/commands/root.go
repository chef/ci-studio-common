package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/cmd/vault-util/internal/pkg/account"
	"github.com/chef/ci-studio-common/internal/pkg/paths"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:          "vault-util",
		Short:        "Utility to access secrets and account information stored in Hashicorp Vault from CI.",
		SilenceUsage: true,
	}

	// These are shared across all the commands
	accountCache = account.NewAccountCache()
)

// Execute handles the execution of child commands and flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("configuration file (default is %s/vault-util.toml)", paths.SettingsDir))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	account.InitConfig()

	viper.SetConfigName("vault-util")      // name of config file (without extension)
	viper.AddConfigPath(paths.SettingsDir) // adding settings directory as first search path
	viper.AddConfigPath(".")               // adding cwd directory as first search path

	// Set environment variable prefix, eg: VAULT_UTIL_GITHUB_TOKEN_NAME
	viper.SetEnvPrefix("vault_util")

	// Override the default config file if a config file has been passed as a flag
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Read in the config file. If we cant find one, we simply proceed with defaults.
	// The config file isn't "required" per say, and the utility should handle when settings are not provided.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal(err)
		}
	}

	// Override our config with any matching environment variables
	viper.AutomaticEnv()
}
