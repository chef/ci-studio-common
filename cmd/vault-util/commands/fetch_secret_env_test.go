package commands

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/MakeNowJust/heredoc"
	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	minify "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

func TestFetchSecretEnvE(t *testing.T) {
	var err error

	m := minify.New()
	m.AddFunc("json", json.Minify)

	fs = filesystem.NewMemFs()
	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	defaultFetchSecretEnvOpts := &fetchSecretEnvOptions{
		format: "sh",
	}

	t.Run("no secrets specified", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Empty(t, output)

		output.Reset()
	})

	t.Run("invalid JSON", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		// JSON has extra `,`
		condensedJSON := `{"VALUE_ONLY":{"value":"true",}}`

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to unmarshall secrets JSON", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		output.Reset()
	})

	t.Run("value was specified", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		vaultUtilJSON := heredoc.Docf(`
			{
				"VALUE_ONLY": {
					"value": "true"
				}
			}
		`)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		expectedOutput := heredoc.Doc(`
			export VALUE_ONLY="true"
		`)
		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output.String())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		output.Reset()
	})

	t.Run("no field parameter was specified", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		vaultUtilJSON := heredoc.Docf(`
			{
				"NO_FIELD": {
					"path": "account/secret/foo"
				}
			}
		`)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "missing required 'field' setting", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		output.Reset()
	})

	t.Run("neither path, value, or account was specified", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		vaultUtilJSON := heredoc.Docf(`
			{
				"NO_PARAM": {
					"field": "foo"
				}
			}
		`)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "insufficient parameters were provided to fetch secret from Vault", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		output.Reset()
	})

	t.Run("complex shell output", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts

		vaultUtilJSON := heredoc.Docf(`
			{
				"VALUE_ONLY": {
					"value": "true"
				},
				"SECRET_PATH": {
					"path": "/my/path",
					"field": "foo"
				},
				"AWS_ACCESS_KEY_ID": {
					"account": "aws/my-account",
					"field": "access_key_id"
				},
				"AWS_SECRET_ACCESS_KEY": {
					"account": "aws/my-account",
					"field": "secret_access_key"
				},
			}
		`)

		specialQuotedSecret := `{"foo":"bar"}`
		accessKeyID := "fake-access-key-id"
		secretAccessKey := "fake-secret-access-key+"

		// Fake Secret
		fakeSecretData := make(map[string]interface{})
		fakeSecretData["foo"] = specialQuotedSecret

		fakeSecret := &vault.Secret{
			Data: fakeSecretData,
		}

		// Fake Account
		fakeAccountData := make(map[string]string)
		fakeAccountData["access_key_id"] = accessKeyID
		fakeAccountData["secret_access_key"] = secretAccessKey

		fakeAccount := &secrets.Account{
			ID:   "aws/my-account",
			Data: fakeAccountData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  fakeSecret,
			err:     nil,
		}

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Contains(t, output.String(), `export VALUE_ONLY="true"`)
		assert.Contains(t, output.String(), fmt.Sprintf("export SECRET_PATH=%q", specialQuotedSecret))
		assert.Contains(t, output.String(), fmt.Sprintf("export AWS_ACCESS_KEY_ID=%q", accessKeyID))
		assert.Contains(t, output.String(), fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=%q", secretAccessKey))
		assert.Equal(t, accountCache["aws/my-account"], fakeAccount)

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("complex powershell output", func(t *testing.T) {
		vaultUtilJSON := heredoc.Docf(`
			{
				"VALUE_ONLY": {
					"value": "true"
				},
				"SECRET_PATH": {
					"path": "/my/path",
					"field": "foo"
				},
				"AWS_ACCESS_KEY_ID": {
					"account": "aws/my-account",
					"field": "access_key_id"
				},
				"AWS_SECRET_ACCESS_KEY": {
					"account": "aws/my-account",
					"field": "secret_access_key"
				},
			}
		`)

		specialQuotedSecret := `{"foo":"bar"}`
		accessKeyID := "fake-access-key-id"
		secretAccessKey := "fake-secret-access-key+"

		// Fake Secret
		fakeSecretData := make(map[string]interface{})
		fakeSecretData["foo"] = specialQuotedSecret

		fakeSecret := &vault.Secret{
			Data: fakeSecretData,
		}

		// Fake Account
		fakeAccountData := make(map[string]string)
		fakeAccountData["access_key_id"] = accessKeyID
		fakeAccountData["secret_access_key"] = secretAccessKey

		fakeAccount := &secrets.Account{
			ID:   "aws/my-account",
			Data: fakeAccountData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  fakeSecret,
			err:     nil,
		}

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		fetchSecretEnvOpts = &fetchSecretEnvOptions{
			format: "ps1",
		}

		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Contains(t, output.String(), `$env:VALUE_ONLY="true"`)
		assert.Contains(t, output.String(), fmt.Sprintf("$env:SECRET_PATH=%q", specialQuotedSecret))
		assert.Contains(t, output.String(), fmt.Sprintf("$env:AWS_ACCESS_KEY_ID=%q", accessKeyID))
		assert.Contains(t, output.String(), fmt.Sprintf("$env:AWS_SECRET_ACCESS_KEY=%q", secretAccessKey))
		assert.Equal(t, accountCache["aws/my-account"], fakeAccount)

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("complex batch output", func(t *testing.T) {
		vaultUtilJSON := heredoc.Docf(`
			{
				"VALUE_ONLY": {
					"value": "true"
				},
				"SECRET_PATH": {
					"path": "/my/path",
					"field": "foo"
				},
				"AWS_ACCESS_KEY_ID": {
					"account": "aws/my-account",
					"field": "access_key_id"
				},
				"AWS_SECRET_ACCESS_KEY": {
					"account": "aws/my-account",
					"field": "secret_access_key"
				},
			}
		`)

		specialQuotedSecret := `{"foo":"bar"}`
		accessKeyID := "fake-access-key-id"
		secretAccessKey := "fake-secret-access-key+"

		// Fake Secret
		fakeSecretData := make(map[string]interface{})
		fakeSecretData["foo"] = specialQuotedSecret

		fakeSecret := &vault.Secret{
			Data: fakeSecretData,
		}

		// Fake Account
		fakeAccountData := make(map[string]string)
		fakeAccountData["access_key_id"] = accessKeyID
		fakeAccountData["secret_access_key"] = secretAccessKey

		fakeAccount := &secrets.Account{
			ID:   "aws/my-account",
			Data: fakeAccountData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  fakeSecret,
			err:     nil,
		}

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		fetchSecretEnvOpts = &fetchSecretEnvOptions{
			format: "batch",
		}

		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Contains(t, output.String(), `"VALUE_ONLY=true"`)
		assert.Contains(t, output.String(), fmt.Sprintf(`"SECRET_PATH=%s"`, specialQuotedSecret))
		assert.Contains(t, output.String(), fmt.Sprintf(`"AWS_ACCESS_KEY_ID=%s"`, accessKeyID))
		assert.Contains(t, output.String(), fmt.Sprintf(`"AWS_SECRET_ACCESS_KEY=%s"`, secretAccessKey))
		assert.Equal(t, accountCache["aws/my-account"], fakeAccount)

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("path secret fetch error", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts
		vaultUtilJSON := heredoc.Docf(`
			{
				"PATH_ONLY": {
					"path": "account/static/foobar",
					"field": "foo"
				}
			}
		`)

		secretsClient = &fakeClient{
			account: nil,
			secret:  nil,
			err:     errors.New("random error"),
		}

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "random error", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		output.Reset()
	})

	t.Run("adds account to cache if doesn't exist", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts
		vaultUtilJSON := heredoc.Docf(`
			{
				"UNCACHED_ACCOUNT": {
					"account": "github/baxterthehacker",
					"field": "token"
				}
			}
		`)

		fakeData := make(map[string]string)
		fakeData["token"] = "password"

		fakeAccount := &secrets.Account{
			ID:   "github/baxterthehacker",
			Data: fakeData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  nil,
			err:     nil,
		}

		assert.Empty(t, accountCache)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		expectedOutput := heredoc.Doc(`
			export UNCACHED_ACCOUNT="password"
		`)
		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output.String())
		assert.Equal(t, fakeAccount, accountCache["github/baxterthehacker"])

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("does not fetch account if exists in cache", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts
		vaultUtilJSON := heredoc.Docf(`
			{
				"CACHED_ACCOUNT": {
					"account": "github/baxterthehacker",
					"field": "token"
				}
			}
		`)

		// Preseed the cache
		cachedData := make(map[string]string)
		cachedData["token"] = "password"

		cachedAccount := &secrets.Account{
			ID:   "github/baxterthehacker",
			Data: cachedData,
		}

		accountCache["github/baxterthehacker"] = cachedAccount

		// Set up a fake to test against
		fakeData := make(map[string]string)
		fakeData["token"] = "bad-password"

		fakeAccount := &secrets.Account{
			ID:   "github/baxterthehacker",
			Data: fakeData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  nil,
			err:     nil,
		}

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		expectedOutput := heredoc.Doc(`
			export CACHED_ACCOUNT="password"
		`)
		err = fetchSecretEnvE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, expectedOutput, output.String())
		assert.Equal(t, cachedAccount, accountCache["github/baxterthehacker"])

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("account fetch error", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts
		vaultUtilJSON := heredoc.Docf(`
			{
				"ACCOUNT_FETCH_ERROR": {
					"account": "github/baxterthehacker",
					"field": "token"
				}
			}
		`)

		secretsClient = &fakeClient{
			account: nil,
			secret:  nil,
			err:     errors.New("random error"),
		}

		assert.Empty(t, accountCache)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "random error", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})

	t.Run("error fetching field from account", func(t *testing.T) {
		fetchSecretEnvOpts = defaultFetchSecretEnvOpts
		vaultUtilJSON := heredoc.Docf(`
			{
				"ACCOUNT_FIELD_FETCH_ERROR": {
					"account": "github/baxterthehacker",
					"field": "token"
				}
			}
		`)

		fakeData := make(map[string]string)
		fakeData["other-token"] = "password"

		fakeAccount := &secrets.Account{
			ID:   "github/baxterthehacker",
			Data: fakeData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  nil,
			err:     nil,
		}

		assert.Empty(t, accountCache)

		condensedJSON, err := m.String("json", vaultUtilJSON)
		assert.Nil(t, err)

		err = os.Setenv("VAULT_UTIL_SECRETS", condensedJSON)
		assert.Nil(t, err)

		err = fetchSecretEnvE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "field 'token' does not exist for the 'github/baxterthehacker' account", err.Error())

		os.Unsetenv("VAULT_UTIL_SECRETS")
		accountCache = make(map[string]*secrets.Account)
		output.Reset()
	})
}
