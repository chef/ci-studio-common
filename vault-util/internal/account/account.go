package account

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type (
	// Cache ensure that we only fetch each account type once
	Cache struct {
		AWSAccounts    map[string]*AWSAccount
		AzureAccounts  map[string]*AzureAccount
		GithubAccounts map[string]*GithubAccount
		GoogleAccounts map[string]*GoogleAccount
	}
)

// NewAccountCache returns a new Cache struct
func NewAccountCache() *Cache {
	return &Cache{
		AWSAccounts:    make(map[string]*AWSAccount),
		AzureAccounts:  make(map[string]*AzureAccount),
		GithubAccounts: make(map[string]*GithubAccount),
		GoogleAccounts: make(map[string]*GoogleAccount),
	}
}

// InitConfig establishes default settings for all acounts
func InitConfig() {
	InitAWSConfig()
	InitAzureConfig()
	InitGithubConfig()
	InitGoogleConfig()
}

// Configure runs commands to configure an account on an instance
func (c *Cache) Configure(accountType string, accountName string) error {
	switch accountType {
	case "aws":
		return c.ConfigureAWS(accountName)
	case "azure":
		return c.ConfigureAzure(accountName)
	case "github":
		return c.ConfigureGithub(accountName)
	case "google":
		return c.ConfigureGoogle(accountName)
	default:
		return fmt.Errorf("unsupported account type: %s", accountType)
	}
}

// GetSecret fetchs a secret from an account in Vault
func (c *Cache) GetSecret(accountVal string, field string) (string, error) {
	bits := strings.Split(accountVal, "/")

	accountName, err := calcAccountName(bits)
	if err != nil {
		return "", err
	}

	switch accountType := bits[0]; accountType {
	case "aws":
		return c.GetAWSSecret(accountName, field)
	case "azure":
		return c.GetAzureSecret(accountName, field)
	case "github":
		return c.GetGithubSecret(accountName, field)
	case "google":
		return c.GetGoogleSecret(accountName, field)
	default:
		return "", fmt.Errorf("unsupported account type: %s", accountType)
	}
}

func calcAccountName(bits []string) (string, error) {
	if len(bits) > 1 {
		return bits[1], nil
	}

	defaultAccountSetting := fmt.Sprintf("%s.default_account", bits[0])

	if !viper.IsSet(defaultAccountSetting) {
		return "", fmt.Errorf("no account name provided, and no default account set for %s", bits[0])
	}

	return viper.GetString(defaultAccountSetting), nil
}
