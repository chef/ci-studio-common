package secrets

import (
	"strings"

	vault "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// VaultClient struct defines a new vault client.
type VaultClient struct {
	client *vault.Client
}

// Client struct defines what the client has access to.
type Client interface {
	GetAccount(id string) (*Account, error)
	GetSecret(path string) (*vault.Secret, error)
}

// NewClient will attempt create a new secrets client.
func NewClient() (*VaultClient, error) {
	client, err := vault.NewClient(nil)
	if err != nil {
		return nil, errors.Wrap(err, "create secrets client")
	}

	return &VaultClient{client: client}, nil
}

// GetAccount will determine which account to use from the specified id.
func (c *VaultClient) GetAccount(id string) (*Account, error) {
	length := 2
	bits := strings.Split(id, "/")

	if len(bits) < length {
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

// GetSecret will read the secret from the given path.
func (c *VaultClient) GetSecret(path string) (*vault.Secret, error) {
	return c.client.Logical().Read(path)
}
