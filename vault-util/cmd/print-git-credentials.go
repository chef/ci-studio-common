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
		Use:   "print-git-credentials [USER]",
		Short: "Utility that will print credentials for a user from Vault in git-credential-helper format.",
		Run:   printCredentials,
	}
)

func init() {
	rootCmd.AddCommand(gitCredentialCmd)
}

func printCredentials(cmd *cobra.Command, args []string) {
	accountName := viper.GetString("github.default_account")
	if len(args) >= 1 {
		accountName = args[0]
	}

	acct, err := account.NewGithubAccount(accountName)
	lib.Check(err)

	fmt.Printf("protocol=https\nhost=github.com\nusername=%s\npassword=%s\n", acct.Name, acct.Token)
}
