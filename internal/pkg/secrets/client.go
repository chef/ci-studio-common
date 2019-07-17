package secrets

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

type VaultClient struct {
	client *vault.Client
}

type Client interface {
	GetAccount(id string) (*Account, error)
	GetSecret(path string) (*vault.Secret, error)
}

func NewClient() (*VaultClient, error) {
	client, err := vault.NewClient(nil)
	if err != nil {
		return nil, err
	}

	return &VaultClient{client: client}, nil
}

func (c *VaultClient) GetAccount(id string) (*Account, error) {
	bits := strings.Split(id, "/")

	if len(bits) < 2 {
		return nil, errors.Errorf("%s is not a valid account specifier", id)
	}

	accountType := bits[0]
	accountName := bits[1]

	switch accountType {
	case "aws":
		return c.newAwsAccount(accountName)
	case "azure":
		return c.newAzureAccount(accountName)
	case "github":
		return c.newGithubAccount(accountName)
	case "google":
		return c.newGoogleAccount(accountName)
	default:
		return nil, errors.Errorf("unsupported account type: %s", accountType)
	}
}

func (c *VaultClient) GetSecret(path string) (*vault.Secret, error) {
	return c.client.Logical().Read(path)
}
