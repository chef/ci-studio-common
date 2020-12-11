package secrets

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Account maps what each account can have.
type Account struct {
	ID   string
	Data map[string]string
}

// Get returns the field within the specified account.
func (a *Account) Get(field string) (string, error) {
	if value, ok := a.Data[field]; ok {
		return value, nil
	}

	return "", errors.Errorf("field '%s' does not exist for the '%s' account", field, a.ID)
}

func (c *VaultClient) newAwsAccount(name string) (*Account, error) {
	secret, err := c.GetSecret(fmt.Sprintf("%s/aws/%s/sts/default", viper.GetString("vault.dynamic_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get aws sts secret")
	}

	data := make(map[string]string)
	data["access_key_id"] = secret.Data["access_key"].(string)
	data["secret_access_key"] = secret.Data["secret_key"].(string)
	data["session_token"] = secret.Data["security_token"].(string)

	return &Account{
		ID:   "aws/" + name,
		Data: data,
	}, nil
}

func (c *VaultClient) newAzureAccount(name string) (*Account, error) {
	secret, err := c.GetSecret(fmt.Sprintf("%s/azure/%s/creds/default", viper.GetString("vault.dynamic_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get azure dynamic secret")
	}

	mappings, err := c.GetSecret(fmt.Sprintf("%s/azure/%s", viper.GetString("vault.static_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get azure static secret")
	}

	data := make(map[string]string)
	data["client_id"] = secret.Data["client_id"].(string)
	data["client_secret"] = secret.Data["client_secret"].(string)
	data["subscription_id"] = mappings.Data["subscription_id"].(string)
	data["tenant_id"] = mappings.Data["tenant_id"].(string)

	return &Account{
		ID:   "azure/" + name,
		Data: data,
	}, nil
}

func (c *VaultClient) newGithubAccount(name string) (*Account, error) {
	secret, err := c.GetSecret(fmt.Sprintf("%s/github/%s", viper.GetString("vault.dynamic_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get github dynamic secret")
	}

	data := make(map[string]string)
	data["token"] = secret.Data["token"].(string)

	return &Account{
		ID:   "github/" + name,
		Data: data,
	}, nil
}

func (c *VaultClient) newGoogleAccount(name string) (*Account, error) {
	base64JSON, err := c.GetSecret(fmt.Sprintf("%s/gcp/%s/key/service-account", viper.GetString("vault.dynamic_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get google dynamic secret")
	}

	decodedJSON, err := base64.StdEncoding.DecodeString(base64JSON.Data["private_key_data"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "google account decode string")
	}

	token, err := c.GetSecret(fmt.Sprintf("%s/gcp/%s/token/access-token", viper.GetString("vault.dynamic_mount"), name))
	if err != nil {
		return nil, errors.Wrap(err, "get google dynamic token")
	}

	data := make(map[string]string)
	data["json"] = string(decodedJSON)
	data["json_base64"] = base64JSON.Data["private_key_data"].(string)
	data["token"] = token.Data["token"].(string)

	return &Account{
		ID:   "google/" + name,
		Data: data,
	}, nil
}
