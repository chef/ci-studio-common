package account

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

	if len(bits) < 2 {
		return "", errors.Errorf("%s is not a valid account specifier", accountVal)
	}

	accountType := bits[0]
	accountName := bits[1]

	switch accountType {
	case "aws":
		return c.GetAWSSecret(accountName, field)
	case "azure":
		return c.GetAzureSecret(accountName, field)
	case "github":
		return c.GetGithubSecret(accountName, field)
	case "google":
		return c.GetGoogleSecret(accountName, field)
	default:
		return "", errors.Errorf("unsupported account type: %s", accountType)
	}
}
