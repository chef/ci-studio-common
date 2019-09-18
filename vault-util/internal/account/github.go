package account

import (
	"fmt"
	"reflect"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/vault"
	"github.com/spf13/viper"
)

// GithubAccount contains patterns and information regarding a GitHub account
type GithubAccount struct {
	Name  string
	Token string
}

var (
	// GithubUserFieldMap maps snake case fields to camelcase struct keys
	GithubUserFieldMap = map[string]string{
		"token": "Token",
	}
)

// InitGithubConfig is called through the chain to establish global defaults
func InitGithubConfig() {
	viper.SetDefault("github.dynamic_mount_root", "account/dynamic/github")
}

// NewGithubAccount creates new instance of GithubAccount
func NewGithubAccount(name string) (*GithubAccount, error) {
	vaultClient, err := vault.NewClient()
	if err != nil {
		return nil, err
	}

	secret, err := vaultClient.Read(fmt.Sprintf("%s/%s/token", viper.GetString("github.dynamic_mount_root"), name))
	if err != nil {
		return nil, err
	}

	return &GithubAccount{Name: name, Token: secret.Data["token"].(string)}, nil
}

func (c *Cache) fetchOrInitGithubAccount(name string) (*GithubAccount, error) {
	var acct *GithubAccount
	var cached bool
	var err error

	if acct, cached = c.GithubAccounts[name]; !cached {
		acct, err = NewGithubAccount(name)

		if err != nil {
			return nil, err
		}

		c.GithubAccounts[name] = acct
	}

	return acct, nil
}

// GetGithubSecret provides a mechanism to quickly/easily get a field from the object
func (c *Cache) GetGithubSecret(name string, field string) (string, error) {
	acct, err := c.fetchOrInitGithubAccount(name)
	if err != nil {
		return "", err
	}

	r := reflect.ValueOf(acct)
	camelField := GithubUserFieldMap[field]
	value := reflect.Indirect(r).FieldByName(camelField)

	return value.String(), nil
}

// ConfigureGithub runs commands to configure the given account using the credential helper API
func (c *Cache) ConfigureGithub(name string) error {
	acct, err := c.fetchOrInitGithubAccount(name)
	if err != nil {
		return err
	}

	cmd := lib.ShellOut("git", "config", "--global", "credential.helper", fmt.Sprintf("!vault-util print-git-credentials %s", acct.Name)).Run()

	if cmd == nil {
		fmt.Printf("Configured git credential helper with '%s' user\n", acct.Name)
	}

	return cmd
}
