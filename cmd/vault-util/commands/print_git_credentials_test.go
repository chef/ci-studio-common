package commands

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/chef/ci-studio-common/internal/pkg/filesystem"
	"github.com/chef/ci-studio-common/internal/pkg/install"
	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

func TestPrintGitCredentailsE(t *testing.T) {
	var err error

	username := "baxterthehacker"
	password := "some-random-token"

	fs = filesystem.NewMemFs()
	ciutils = install.DefaultInstall()
	filename := ciutils.SettingsPath("github_token_baxterthehacker")

	cmd := &cobra.Command{}
	args := []string{}
	output := new(bytes.Buffer)
	cmd.SetOutput(output)

	printGitCredentialsOpts = &printGitCredentialsOptions{
		account: username,
	}

	t.Run("cache file does not exist", func(t *testing.T) {
		fakeData := make(map[string]string)
		fakeData["token"] = password

		fakeAccount := &secrets.Account{
			ID:   "github/" + username,
			Data: fakeData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  nil,
			err:     nil,
		}

		err = printGitCredentialsE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("username=x-access-token\npassword=%s", password), output.String())

		contents, err := fs.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, password, strings.TrimSpace(string(contents)))

		output.Reset()
		err = fs.Remove(filename)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache file exists", func(t *testing.T) {
		err := fs.WriteFile(filename, []byte(password), 0600)
		if err != nil {
			t.Fatal(err)
		}

		err = printGitCredentialsE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("username=x-access-token\npassword=%s", password), output.String())

		contents, err := fs.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, password, strings.TrimSpace(string(contents)))

		output.Reset()
		err = fs.Remove(filename)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("cache file exists but is expired", func(t *testing.T) {
		err := fs.WriteFile(filename, []byte("old-password"), 0600)
		if err != nil {
			t.Fatal(err)
		}

		oneHourAgo := time.Now().Add(-60 * time.Minute)
		err = fs.Chtimes(filename, oneHourAgo, oneHourAgo)
		if err != nil {
			t.Fatal(err)
		}

		fakeData := make(map[string]string)
		fakeData["token"] = password

		fakeAccount := &secrets.Account{
			ID:   "github/" + username,
			Data: fakeData,
		}

		secretsClient = &fakeClient{
			account: fakeAccount,
			secret:  nil,
			err:     nil,
		}

		err = printGitCredentialsE(cmd, args)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("username=x-access-token\npassword=%s", password), output.String())

		contents, err := fs.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, password, strings.TrimSpace(string(contents)))

		output.Reset()
		err = fs.Remove(filename)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("account not found", func(t *testing.T) {
		secretsClient = &fakeClient{
			account: nil,
			secret:  nil,
			err:     errors.New("some internal error"),
		}

		err = printGitCredentialsE(cmd, args)
		assert.NotNil(t, err)

		output.Reset()
	})
}
