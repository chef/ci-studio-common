package vault

import (
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

var vaultConfig = &api.Config{
	Address: os.Getenv("VAULT_ADDR"),
}

// Client is a wrapper around the vault.Client object
type Client struct {
	vault *api.Client
}

// NewClient returns a new Vault client
//
// For right now, we assume that you have the token stored someplace the API can access it.
// In the future we'll want to support fetching the token from instance profile, but our
// current use case doesn't need to support that.
func NewClient() (*Client, error) {
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	if client.Token() == "" {
		return nil, errors.New("the VAULT_TOKEN environment variable must be set")
	}

	return &Client{vault: client}, nil
}

// Read provides a short cut to the logical read of the native Vault client
func (client *Client) Read(path string) (*api.Secret, error) {
	return client.vault.Logical().Read(path)
}
