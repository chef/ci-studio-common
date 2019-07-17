package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
)

var (
	configureAccountsCmd = &cobra.Command{
		Use:   "configure-accounts",
		Short: "Configure the accounts specified in the VAULT_UTIL_ACCOUNTS environment variable.",
		RunE:  configureAccountsE,
	}
)

func init() {
	rootCmd.AddCommand(configureAccountsCmd)
}

// Some utilities do not like it when you try and configure accounts at the same time.
// If we're able to get the first lock on the instance, proceed.
// Otherwise, we'll loop and wait for the update to finish (lock is released)
func configureAccountsE(cmd *cobra.Command, args []string) error {
	var accountsJSON map[string][]string

	accountsEnv := os.Getenv("VAULT_UTIL_ACCOUNTS")

	if accountsEnv == "" {
		return nil
	}

	lock, err := fslock.GetLock(ciutils.LockPath("configure-accounts"))
	if err != nil {
		return err
	}

	retryErr := retry.Do(
		func() error {
			_, lockErr := lock.TryLock()
			if lockErr != nil {
				return lockErr
			}

			jsonErr := json.Unmarshal([]byte(accountsEnv), &accountsJSON)
			if jsonErr != nil {
				return errors.Wrap(jsonErr, "failed to unmarshal accounts JSON")
			}

			configErr := configureAccounts(accountsJSON)
			if configErr != nil {
				return errors.Wrap(configErr, "failed to configure accounts")
			}

			unlockErr := lock.Unlock()
			if unlockErr != nil {
				return errors.Wrap(unlockErr, "failed to release account configuration lock")
			}

			return nil
		},
		retry.Attempts(fslock.GetRetryAttempts()),
		retry.Delay(fslock.GetRetryDelay()),
		retry.DelayType(fslock.GetRetryDelayType()),
		retry.RetryIf(func(err error) bool {
			if err == filesystem.ErrLocked {
				return true
			}
			return false
		}),
		retry.OnRetry(func(n uint, err error) {
			if err == filesystem.ErrLocked {
				cmd.Printf("another account configuration already in progress -- waiting (%d/%d)\n", n, fslock.GetRetryAttempts())
			}
		}),
	)

	return retryErr
}

func configureAccounts(accountsJSON map[string][]string) error {
	for accountType, accountVals := range accountsJSON {
		for _, accountVal := range accountVals {
			switch accountType {
			case "aws":
				return configureAws(accountVal)
			case "github":
				return configureGithub(accountVal)
			default:
				return errors.Errorf("unsupported account type: %s", accountType)
			}
		}
	}

	return nil
}

func configureAws(name string) error {
	account, err := secretsClient.GetAccount("aws/" + name)
	if err != nil {
		return err
	}

	if err = execCommand("aws", "configure", "set", "aws_access_key_id", account.Data["access_key_id"], "--profile", name).Run(); err != nil {
		return err
	}

	if err = execCommand("aws", "configure", "set", "aws_secret_access_key", account.Data["secret_access_key"], "--profile", name).Run(); err != nil {
		return err
	}

	if err = execCommand("aws", "configure", "set", "aws_session_token", account.Data["session_token"], "--profile", name).Run(); err != nil {
		return err
	}

	return execCommand("aws", "configure", "set", "region", viper.GetString("aws.region"), "--profile", name).Run()
}

func configureGithub(name string) error {
	return execCommand("git", "config", "--global", "credential.helper", fmt.Sprintf("!vault-util print-git-credentials --app %s", name)).Run()
}
