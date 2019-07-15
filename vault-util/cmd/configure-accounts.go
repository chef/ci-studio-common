package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/avast/retry-go"
	"github.com/chef/ci-studio-common/lib"
	"github.com/juju/fslock"
	"github.com/spf13/cobra"
)

var (
	configureJobCmd = &cobra.Command{
		Use:   "configure-accounts",
		Short: "Configure the accounts specified in the VAULT_UTIL_ACCOUNTS environment variable.",
		Run:   tryToConfigureAccounts,
	}
)

func init() {
	rootCmd.AddCommand(configureJobCmd)
}

// Some utilities do not like it when you try and configure accounts at the same time.
// If we're able to get the first lock on the instance, proceed.
// Otherwise, we'll loop and wait for the update to finish (lock is released)
func tryToConfigureAccounts(cmd *cobra.Command, args []string) {
	accountsEnv := os.Getenv("VAULT_UTIL_ACCOUNTS")

	if accountsEnv == "" {
		return
	}

	lock := fslock.New(lib.LockPath("configure-accounts"))

	err := retry.Do(
		func() error {
			lockErr := lock.TryLock()

			if lockErr == nil {
				configureAccounts(accountsEnv)

				unlockErr := lock.Unlock()
				lib.Check(unlockErr)

				return nil
			}

			fmt.Println("another account configuration already in progress -- waiting")
			return lockErr
		},
	)
	lib.Check(err)
}

func configureAccounts(accountsEnv string) {
	var accountsJSON map[string]interface{}

	err := json.Unmarshal([]byte(accountsEnv), &accountsJSON)
	lib.Check(err)

	for accountType, accountVals := range accountsJSON {
		for _, accountVal := range accountVals.([]interface{}) {
			err := accountCache.Configure(accountType, accountVal.(string))
			lib.Check(err)
		}
	}
}
