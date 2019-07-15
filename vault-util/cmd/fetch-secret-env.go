package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/vault"
)

var (
	fetchSecretsCmd = &cobra.Command{
		Use:   "fetch-secret-env",
		Short: "Fetch the secrets specified in the VAULT_UTIL_SECRETS environment variable from Vault.",
		Run:   printSecrets,
	}

	fetchSecretsOpts = struct {
		format string
	}{}

	vaultClient *vault.Client
	secretsJSON map[string]interface{}
)

func init() {
	var err error

	vaultClient, err = vault.NewClient()
	lib.Check(err)

	rootCmd.AddCommand(fetchSecretsCmd)

	fetchSecretsCmd.Flags().StringVar(&fetchSecretsOpts.format, "format", "sh", "Which format to use when exporting environment variables.")
}

func printSecrets(cmd *cobra.Command, args []string) {
	secretsEnv := os.Getenv("VAULT_UTIL_SECRETS")

	if len(secretsEnv) == 0 {
		return
	}

	err := json.Unmarshal([]byte(secretsEnv), &secretsJSON)
	lib.Check(err)

	for k, v := range secretsJSON {
		value := fetchSecret(v.(map[string]interface{}))

		switch fetchSecretsOpts.format {
		case "sh":
			fmt.Printf("export %s=%q\n", k, value)
		case "ps1":
			fmt.Printf("$env:%s=%q\n", k, value)
		case "batch":
			fmt.Printf("\"%s=%s\"\n", k, value)
		default:
			log.Fatalf("'%s' is not a supported format", fetchSecretsOpts.format)
		}
	}
}

func fetchSecret(settings map[string]interface{}) string {
	if value, exists := settings["value"]; exists {
		return value.(string)
	}

	field, exists := settings["field"]

	if !exists {
		log.Fatal(errors.New("missing required 'field' setting"))

	}

	if path, exists := settings["path"]; exists {
		secret, err := vaultClient.Read(path.(string))
		lib.Check(err)

		return secret.Data[field.(string)].(string)
	}

	if accountVal, exists := settings["account"]; exists {
		secret, err := accountCache.GetSecret(accountVal.(string), field.(string))
		lib.Check(err)

		return secret
	}

	log.Fatal(errors.New("insufficient parameters were provided to fetch secret from Vault"))

	return ""
}
