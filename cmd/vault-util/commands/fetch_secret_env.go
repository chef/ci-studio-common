package commands

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

type (
	fetchSecretEnvOptions struct {
		format string
	}

	Secret struct {
		Field   string `json:"field"`
		Path    string `json:"path"`
		Value   string `json:"value"`
		Account string `json:"account"`
	}
)

var (
	accountCache = make(map[string]*secrets.Account)

	fetchSecretEnvCmd = &cobra.Command{
		Use:   "fetch-secret-env",
		Short: "Fetch the secrets specified in the VAULT_UTIL_SECRETS environment variable from Vault.",
		RunE:  fetchSecretEnvE,
	}

	fetchSecretEnvOpts = &fetchSecretEnvOptions{}
)

func init() {
	rootCmd.AddCommand(fetchSecretEnvCmd)

	fetchSecretEnvCmd.Flags().StringVar(&fetchSecretEnvOpts.format, "format", "sh", "Which format to use when exporting environment variables.")
}

func fetchSecretEnvE(cmd *cobra.Command, args []string) error {
	var secretsJSON map[string]Secret

	secretsEnv := os.Getenv("VAULT_UTIL_SECRETS")

	if len(secretsEnv) == 0 {
		return nil
	}

	err := json.Unmarshal([]byte(secretsEnv), &secretsJSON)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall secrets JSON")
	}

	for k, v := range secretsJSON {
		value, err := fetchSecret(v)
		if err != nil {
			return errors.Wrap(err, "failed to print out secrets")
		}

		switch fetchSecretEnvOpts.format {
		case "sh":
			cmd.Printf("export %s=%q\n", k, value)
		case "ps1":
			cmd.Printf("$env:%s=%q\n", k, value)
		case "batch":
			cmd.Printf("\"%s=%s\"\n", k, value)
		default:
			return errors.Errorf("'%s' is not a supported format", fetchSecretEnvOpts.format)
		}
	}

	return nil
}

func fetchSecret(s Secret) (string, error) {
	var err error

	if s.Value != "" {
		return s.Value, nil
	}

	if s.Field == "" {
		return "", errors.New("missing required 'field' setting")
	}

	if s.Path != "" {
		secret, err := secretsClient.GetSecret(s.Path)
		if err != nil {
			return "", err
		}

		return secret.Data[s.Field].(string), nil
	}

	if s.Account != "" {
		var account *secrets.Account
		var exists bool

		account, exists = accountCache[s.Account]

		// We do not want to re-fetch an account everytime we load a secret.
		// This can cause issues because things aws access_keys and secret_keys wouldn't match.
		// To avoid this, we cache accounts. Only fetch an account if its not in the cache.
		if !exists {
			account, err = secretsClient.GetAccount(s.Account)
			if err != nil {
				return "", err
			}

			accountCache[s.Account] = account
		}

		secret, err := account.Get(s.Field)
		if err != nil {
			return "", err
		}

		return secret, nil
	}

	return "", errors.New("insufficient parameters were provided to fetch secret from Vault")
}
