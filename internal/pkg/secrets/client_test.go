package secrets

import (
	"bytes"
	"net"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()

	output := new(bytes.Buffer)

	secretsClient := &VaultClient{
		client: client,
	}

	t.Run("secret exists", func(t *testing.T) {
		secret, err := secretsClient.GetSecret("secret/test")

		data := make(map[string]string)
		data["key"] = secret.Data["key"].(string)

		assert.Nil(t, err)
		assert.Equal(t, data["key"], "value")

		output.Reset()
	})
}

func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	// Setup required secrets, policies, etc.
	_, err = client.Logical().Write("secret/test", map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}
