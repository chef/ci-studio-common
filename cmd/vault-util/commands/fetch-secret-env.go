package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/cmd/vault-util/internal/pkg/vault"
)

type secret struct {
	Field   string `json:"field"`
	Path    string `json:"path"`
	Value   string `json:"value"`
	Account string `json:"account"`
}

var (
	fetchSecretsCmd = &cobra.Command{
		Use:   "fetch-secret-env",
		Short: "Fetch the secrets specified in the VAULT_UTIL_SECRETS environment variable from Vault.",
		RunE:  printSecrets,
	}

	fetchSecretsOpts = struct {
		format string
	}{}

	secretsJSON map[string]secret
)

func init() {
	rootCmd.AddCommand(fetchSecretsCmd)

	fetchSecretsCmd.Flags().StringVar(&fetchSecretsOpts.format, "format", "sh", "Which format to use when exporting environment variables.")
}

func printSecrets(cmd *cobra.Command, args []string) error {
	secretsEnv := os.Getenv("VAULT_UTIL_SECRETS")

	if len(secretsEnv) == 0 {
		return nil
	}

	vaultClient, err := vault.NewClient()
	if err != nil {
		return errors.Wrap(err, "failed to create vault client")
	}

	err = json.Unmarshal([]byte(secretsEnv), &secretsJSON)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshall secrets JSON")
	}

	for k, v := range secretsJSON {
		value, err := fetchSecret(vaultClient, v)
		if err != nil {
			return errors.Wrap(err, "failed to print out secrets")
		}

		switch fetchSecretsOpts.format {
		case "sh":
			fmt.Printf("export %s=%q\n", k, value)
		case "ps1":
			fmt.Printf("$env:%s=%q\n", k, value)
		case "batch":
			fmt.Printf("\"%s=%s\"\n", k, value)
		default:
			return errors.Errorf("'%s' is not a supported format", fetchSecretsOpts.format)
		}
	}

	return nil
}

func fetchSecret(client *vault.Client, s secret) (string, error) {
	if s.Value != "" {
		return s.Value, nil
	}

	if s.Field == "" {
		return "", errors.New("missing required 'field' setting")
	}

	if s.Path != "" {
		vaultSecret, err := client.Read(s.Path)
		if err != nil {
			return "", errors.Wrapf(err, "failed to read secret from path %s", s.Path)
		}

		return vaultSecret.Data[s.Field].(string), nil
	}

	if s.Account != "" {
		secret, err := accountCache.GetSecret(s.Account, s.Field)
		if err != nil {
			return "", errors.Wrapf(err, "failed to fetch %s account", s.Account)
		}

		return secret, nil
	}

	return "", errors.New("insufficient parameters were provided to fetch secret from Vault")
}
