package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	minify "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

var configureCommands []string

func TestConfigureAccountsE(t *testing.T) {
	var err error

	m := minify.New()
	m.AddFunc("json", json.Minify)

	vaultUtilJSON := heredoc.Docf(`
		{
			"aws": ["my-account"]
		}
	`)
	defaultCondensedJSON, err := m.String("json", vaultUtilJSON)
	assert.Nil(t, err)

	fs = filesystem.NewMemFs()
	ciutils = install.DefaultInstall()

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	t.Run("no accounts specified", func(t *testing.T) {
		err = configureAccountsE(cmd, args)
		assert.Nil(t, err)
		assert.Empty(t, output)

		output.Reset()
	})

	t.Run("lock already exists", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", defaultCondensedJSON)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		lock, _ := fslock.GetLock(ciutils.LockPath("configure-accounts"))
		_, err = lock.TryLock()
		assert.Nil(t, nil)

		expectedOutput := heredoc.Docf(`
			another account configuration already in progress -- waiting (0/3)
			another account configuration already in progress -- waiting (1/3)
			another account configuration already in progress -- waiting (2/3)
		`)
		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Equal(t, expectedOutput, output.String())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("failed to get lock", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", defaultCondensedJSON)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
			Err:            errors.New("internal go-filelock error"),
		}

		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Equal(t, "internal go-filelock error", err.Error())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("fail to unmarshall JSON", func(t *testing.T) {
		// missing a final quote
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"aws":["my-account]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Empty(t, output.String())
		assert.Regexp(t, "failed to unmarshal accounts JSON", err.Error())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("unsupported account", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"foobar":["my-account"]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "unsupported account type", err.Error())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("error configuring aws account", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"aws":["my-account"]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		configureCommands = []string{}
		execCommand = fakeFailureConfigureCommand
		defer func() { execCommand = exec.Command }()

		accessKeyID := "fake-access-key-id"
		secretAccessKey := "fake-secret-access-key+"
		sessionToken := "fake-secret-token"

		fakeAccountData := make(map[string]string)
		fakeAccountData["access_key_id"] = accessKeyID
		fakeAccountData["secret_access_key"] = secretAccessKey
		fakeAccountData["session_token"] = sessionToken

		fakeAccount := &secrets.Account{
			ID:   "aws/my-account",
			Data: fakeAccountData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
		}

		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to configure accounts", err.Error())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("no error configuring aws account", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"aws":["my-account"]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		configureCommands = []string{}
		execCommand = fakeSuccessConfigureCommand
		defer func() { execCommand = exec.Command }()

		accessKeyID := "fake-access-key-id"
		secretAccessKey := "fake-secret-access-key+"
		sessionToken := "fake-secret-token"

		fakeAccountData := make(map[string]string)
		fakeAccountData["access_key_id"] = accessKeyID
		fakeAccountData["secret_access_key"] = secretAccessKey
		fakeAccountData["session_token"] = sessionToken

		fakeAccount := &secrets.Account{
			ID:   "aws/my-account",
			Data: fakeAccountData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
		}

		err = configureAccountsE(cmd, args)
		assert.Nil(t, err)

		assert.Equal(t, "aws configure set aws_access_key_id fake-access-key-id --profile my-account", configureCommands[0])
		assert.Equal(t, "aws configure set aws_secret_access_key fake-secret-access-key+ --profile my-account", configureCommands[1])
		assert.Equal(t, "aws configure set aws_session_token fake-secret-token --profile my-account", configureCommands[2])
		assert.Equal(t, "aws configure set region us-east-1 --profile my-account", configureCommands[3])

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("error configuring git account", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"github":["my-account"]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		configureCommands = []string{}
		execCommand = fakeFailureConfigureCommand
		defer func() { execCommand = exec.Command }()

		err = configureAccountsE(cmd, args)
		assert.NotNil(t, err)
		assert.Regexp(t, "failed to configure accounts", err.Error())

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})

	t.Run("no error configuring git account", func(t *testing.T) {
		err = os.Setenv("VAULT_UTIL_ACCOUNTS", `{"github":["my-app"]}`)
		assert.Nil(t, err)

		fslock = &filesystem.MemMapLock{
			RetryAttempts:  3,
			RetryDelay:     time.Nanosecond,
			RetryDelayType: retry.FixedDelay,
			Locks:          make(map[string]*filesystem.MemLock),
		}

		configureCommands = []string{}
		execCommand = fakeSuccessConfigureCommand
		defer func() { execCommand = exec.Command }()

		err = configureAccountsE(cmd, args)
		assert.Nil(t, err)

		assert.Equal(t, "git config --global credential.helper !vault-util print-git-credentials --app my-app", configureCommands[0])

		os.Unsetenv("VAULT_UTIL_ACCOUNTS")
		output.Reset()
	})
}

func fakeSuccessConfigureCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestConfigureAccountSuccessHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	configureCommands = append(configureCommands, fmt.Sprintf("%s %s", command, strings.Join(args, " ")))

	return cmd
}

func fakeFailureConfigureCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestConfigureAccountFailureHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}

	configureCommands = append(configureCommands, fmt.Sprintf("%s %s", command, strings.Join(args, " ")))

	return cmd
}

func TestConfigureAccountSuccessHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestConfigureAccountFailureHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(1)
}
