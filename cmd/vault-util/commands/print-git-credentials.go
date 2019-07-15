package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/cmd/vault-util/internal/pkg/account"
)

var (
	gitCredentialCmd = &cobra.Command{
		Use:   "print-git-credentials [USER]",
		Short: "Utility that will print credentials for a user from Vault in git-credential-helper format.",
		RunE:  printCredentials,
	}
)

func init() {
	rootCmd.AddCommand(gitCredentialCmd)
}

func printCredentials(cmd *cobra.Command, args []string) error {
	accountName := viper.GetString("github.default_account")
	if len(args) >= 1 {
		accountName = args[0]
	}

	acct, err := account.NewGithubAccount(accountName)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch credentials for %s GitHub account", accountName)
	}

	fmt.Printf("protocol=https\nhost=github.com\nusername=%s\npassword=%s\n", acct.Name, acct.Token)

	return nil
}
