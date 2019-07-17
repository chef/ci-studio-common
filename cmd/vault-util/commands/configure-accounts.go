package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/avast/retry-go"
	"github.com/juju/fslock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/paths"
)

var (
	configureJobCmd = &cobra.Command{
		Use:   "configure-accounts",
		Short: "Configure the accounts specified in the VAULT_UTIL_ACCOUNTS environment variable.",
		RunE:  tryToConfigureAccounts,
	}

	accountsJSON map[string][]string
)

func init() {
	rootCmd.AddCommand(configureJobCmd)
}

// Some utilities do not like it when you try and configure accounts at the same time.
// If we're able to get the first lock on the instance, proceed.
// Otherwise, we'll loop and wait for the update to finish (lock is released)
func tryToConfigureAccounts(cmd *cobra.Command, args []string) error {
	accountsEnv := os.Getenv("VAULT_UTIL_ACCOUNTS")

	if accountsEnv == "" {
		return nil
	}

	lock := fslock.New(paths.LockPath("configure-accounts"))

	err := retry.Do(
		func() error {
			lockErr := lock.TryLock()

			if lockErr == nil {
				configErr := configureAccounts(accountsEnv)
				if configErr != nil {
					return errors.Wrap(configErr, "failed to configure accounts")
				}

				unlockErr := lock.Unlock()
				if unlockErr != nil {
					return errors.Wrap(unlockErr, "failed to release account configuration lock")
				}

				return nil
			}

			fmt.Println("another account configuration already in progress -- waiting")
			return lockErr
		},
	)

	return err
}

func configureAccounts(accountsEnv string) error {
	err := json.Unmarshal([]byte(accountsEnv), &accountsJSON)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal accounts JSON")
	}

	for accountType, accountVals := range accountsJSON {
		for _, accountVal := range accountVals {
			err = accountCache.Configure(accountType, accountVal)
			if err != nil {
				return errors.Wrapf(err, "failed to configure %s %s account", accountType, accountVal)
			}
		}
	}

	return nil
}
