package cmd

import (
	"fmt"
	"os"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/account"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Execute does something?
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("configuration file (default is %s/vault-util.toml)", lib.SettingsDir))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	account.InitConfig()

	viper.SetConfigName("vault-util")      // name of config file (without extension)
	viper.AddConfigPath(lib.SettingsDir) // adding settings directory as first search path
	viper.AddConfigPath(".")               // adding cwd directory as first search path

	// Set environment variable prefix, eg: AUTOMATE_CONFIG_MGMT_PORT
	viper.SetEnvPrefix("ci_studio_common")

	// Override the default config file if a config file has been passed as a flag
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	// Read in the config file. If we cant find one, we simply proceed with defaults.
	// The config file isn't "required" per say, and the utility should handle when settings are not provided.
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// This is "ok" -- do nothing and proceed with defaults
		default:
			panic(err)
		}

	}

	// Override our config with any matching environment variables
	viper.AutomaticEnv()
}
