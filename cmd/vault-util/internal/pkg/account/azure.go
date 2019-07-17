package account

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/chef/ci-studio-common/cmd/vault-util/internal/pkg/vault"
)

// AzureAccount holds the credentials for Azure
type AzureAccount struct {
	Name           string
	ClientID       string
	ClientSecret   string
	SubscriptionID string
	TenantID       string
}

// InitAzureConfig is called through the chain to establish global defaults
func InitAzureConfig() {
	viper.SetDefault("azure.dynamic_mount_root", "account/dynamic/azure")
	viper.SetDefault("azure.static_mount_root", "account/static/azure")
}

// NewAzureAccount returns a new AzureAccount struct
func NewAzureAccount(name string) (*AzureAccount, error) {
	vaultClient, err := vault.NewClient()
	if err != nil {
		return nil, err
	}

	secret, err := vaultClient.Read(fmt.Sprintf("%s/%s/creds/default", viper.GetString("azure.dynamic_mount_root"), name))
	if err != nil {
		return nil, err
	}

	mappings, err := vaultClient.Read(fmt.Sprintf("%s/%s", viper.GetString("azure.static_mount_root"), name))
	if err != nil {
		return nil, err
	}

	return &AzureAccount{
		Name:           name,
		ClientID:       secret.Data["client_id"].(string),
		ClientSecret:   secret.Data["client_secret"].(string),
		SubscriptionID: mappings.Data["subscription_id"].(string),
		TenantID:       mappings.Data["tenant_id"].(string),
	}, nil
}

func (c *Cache) fetchOrInitAzureAccount(name string) (*AzureAccount, error) {
	var acct *AzureAccount
	var cached bool
	var err error

	if acct, cached = c.AzureAccounts[name]; !cached {
		acct, err = NewAzureAccount(name)

		if err != nil {
			return nil, err
		}

		c.AzureAccounts[name] = acct
	}

	return acct, nil
}

// ConfigureAzure runs commands to configure the given account
func (c *Cache) ConfigureAzure(name string) error {
	return fmt.Errorf("azure cli configuration not supported")
}

// GetAzureSecret provides a mechanism to quickly fetch field
func (c *Cache) GetAzureSecret(name string, field string) (string, error) {
	acct, err := c.fetchOrInitAzureAccount(name)
	if err != nil {
		return "", err
	}

	switch field {
	case "client_id":
		return acct.ClientID, nil
	case "client_secret":
		return acct.ClientSecret, nil
	case "subscription_id":
		return acct.SubscriptionID, nil
	case "tenant_id":
		return acct.TenantID, nil
	default:
		return "", errors.Errorf("%s is not a supported field for azure accounts", field)
	}
}
