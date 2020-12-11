package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type printGitCredentialsOptions struct {
	account string
}

var (
	printGitCredentialsCmd = &cobra.Command{
		Use:   "print-git-credentials",
		Short: "Utility that will print credentials for a GitHub App from Vault in git-credential-helper format.",
		RunE:  printGitCredentialsE,
	}

	printGitCredentialsOpts = &printGitCredentialsOptions{}
)

func init() {
	rootCmd.AddCommand(printGitCredentialsCmd)

	printGitCredentialsCmd.Flags().StringVar(&printGitCredentialsOpts.account, "app", "", "GitHub App to use to generate credentials (required)")
	if err := printGitCredentialsCmd.MarkFlagRequired("app"); err != nil {
		log.Fatal(err)
	}
}

func printGitCredentialsE(cmd *cobra.Command, args []string) error {
	filename := ciutils.SettingsPath(fmt.Sprintf("github_token_%s", printGitCredentialsOpts.account))

	token, err := tokenFromCache(filename)
	if err != nil {
		token, err = reloadTokenFromVault(filename)
		if err != nil {
			return err
		}
	}

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
		return "", errors.New("github token has expired")
	}

	tokenRaw, err := fs.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tokenRaw)), nil
}

func reloadTokenFromVault(filename string) (string, error) {
	account, err := secretsClient.GetAccount("github/" + printGitCredentialsOpts.account)
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
