package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:          "get",
	Short:        "Fetch GitHub credentials from GitHub using the GitHub App secrets plugin.",
	SilenceUsage: true,
	RunE:         getE,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func getE(cmd *cobra.Command, args []string) error {
	filename := ciutils.SettingsPath(fmt.Sprintf("github_token_%s", rootOpts.install))

	token, err := tokenFromCache(filename)
	if err != nil {
		token, err = reloadTokenFromVault(filename)
		if err != nil {
			return err
		}
	}

	cmd.SetOutput(os.Stdout)
	cmd.Printf("username=x-access-token\npassword=%s", token)
	return nil
}

func tokenFromCache(filename string) (string, error) {
	info, err := fs.Stat(filename)
	if err != nil {
		return "", err
	}

	// Tokens technically live 60 minutes, but we use 55 for wiggle room
	if time.Since(info.ModTime()) > 55*time.Minute {
		return "", fmt.Errorf("github token has expired")
	}

	tokenRaw, err := fs.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tokenRaw)), nil
}

func reloadTokenFromVault(filename string) (string, error) {
	account, err := secretsClient.GetAccount("github/" + rootOpts.install)
	if err != nil {
		return "", err
	}

	accountToken, err := account.Get("token")
	if err != nil {
		return "", err
	}

	err = fs.WriteFile(filename, []byte(accountToken), 0600)
	if err != nil {
		return "", err
	}

	return accountToken, nil
}
