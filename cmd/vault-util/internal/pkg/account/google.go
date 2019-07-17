package account

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/cmd/vault-util/internal/pkg/vault"
)

// GoogleAccount holds credentials for AWS account
type GoogleAccount struct {
	Name  string
	Token string
	JSON  string
}

// InitGoogleConfig is called through the chain to establish global defaults
func InitGoogleConfig() {
	viper.SetDefault("google.dynamic_mount_root", "account/dynamic/gcp")
}

// NewGoogleAccount returns new GoogleAccount struct
func NewGoogleAccount(name string) (*GoogleAccount, error) {
	vaultClient, err := vault.NewClient()
	if err != nil {
		return nil, err
	}

	base64JSON, err := vaultClient.Read(fmt.Sprintf("%s/%s/key/service-account", viper.GetString("google.dynamic_mount_root"), name))
	if err != nil {
		return nil, err
	}

	decodedJSON, err := base64.StdEncoding.DecodeString(base64JSON.Data["private_key_data"].(string))
	if err != nil {
		return nil, err
	}

	token, err := vaultClient.Read(fmt.Sprintf("%s/%s/token/access-token", viper.GetString("google.dynamic_mount_root"), name))
	if err != nil {
		return nil, err
	}

	return &GoogleAccount{
		Name:  name,
		Token: token.Data["token"].(string),
		JSON:  string(decodedJSON),
	}, nil
}

// ConfigureGoogle runs commands to configure the given account
func (c *Cache) ConfigureGoogle(name string) error {
	return fmt.Errorf("gcp cli configuration not supported")
}

// GetGoogleSecret provides a mechanism to quickly fetch field
func (c *Cache) GetGoogleSecret(name string, field string) (string, error) {
	var acct *GoogleAccount
	var cached bool
	var err error

	if acct, cached = c.GoogleAccounts[name]; !cached {
		acct, err = NewGoogleAccount(name)

		if err != nil {
			return "", err
		}

		c.GoogleAccounts[name] = acct
	}

	switch field {
	case "token":
		return acct.Token, nil
	case "json":
		return acct.JSON, nil
	default:
		return "", errors.Errorf("%s is not a supported field for google accounts", field)
	}
}
