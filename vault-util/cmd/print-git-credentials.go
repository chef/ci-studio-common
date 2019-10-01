package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/account"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/djherbis/times.v1"
)

var (
	gitCredentialCmd = &cobra.Command{
		Use:   "print-git-credentials [APP]",
		Short: "Utility that will print credentials for a GitHub App from Vault in git-credential-helper format.",
		RunE:  printCredentialsE,
	}
)

func init() {
	rootCmd.AddCommand(gitCredentialCmd)
}

func printCredentialsE(cmd *cobra.Command, args []string) error {
	appName := viper.GetString("github.default_app")
	if len(args) == 2 {
		appName = args[0]
	}

	filename := lib.SettingsPath(fmt.Sprintf("github_token_%s", appName))

	var token string

	times, err := times.Stat(filename)
	if err != nil {
		token, err = reloadTokenFromVault(filename, appName)
		if err != nil {
			return err
		}
	} else {
		// The timeout is 1hr, but give ourselves some wiggle room
		someMinutesAgo := time.Now().Add(time.Duration(-55) * time.Minute)

		if times.ChangeTime().Before(someMinutesAgo) {
			token, err = reloadTokenFromVault(filename, appName)
			if err != nil {
				return err
			}
		} else {
			tokenRaw, err := ioutil.ReadFile(filename)
			if err != nil {
				token, err = reloadTokenFromVault(filename, appName)
				if err != nil {
					return err
				}
			} else {
				token = strings.TrimSpace(string(tokenRaw))
			}
		}
	}

	fmt.Printf("username=x-access-token\npassword=%s", token)
	return nil
}

func reloadTokenFromVault(filename string, appName string) (string, error) {
	acct, err := account.NewGithubAccount(appName)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(filename, []byte(acct.Token), 0600)
	if err != nil {
		return "", err
	}

	return acct.Token, nil
}
