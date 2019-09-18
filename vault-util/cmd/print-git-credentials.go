package cmd

import (
	"fmt"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/account"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	gitCredentialCmd = &cobra.Command{
		Use:   "print-git-credentials [APP]",
		Short: "Utility that will print credentials for a GitHub App from Vault in git-credential-helper format.",
		Run:   printCredentials,
	}
)

func init() {
	rootCmd.AddCommand(gitCredentialCmd)
}

func printCredentials(cmd *cobra.Command, args []string) {
	appName := viper.GetString("github.default_app")
	if len(args) >= 1 {
		appName = args[0]
	}

	acct, err := account.NewGithubAccount(appName)
	lib.Check(err)

	fmt.Printf("protocol=https\nhost=github.com\nusername=x-access-token\npassword=%s\n", acct.Token)
}
