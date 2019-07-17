package commands

import (
	vault "github.com/hashicorp/vault/api"

	"github.com/chef/ci-studio-common/internal/pkg/secrets"
)

type fakeClient struct {
	account *secrets.Account
	secret  *vault.Secret
	err     error
}

func (f *fakeClient) GetAccount(id string) (*secrets.Account, error) {
	return f.account, f.err
}

func (f *fakeClient) GetSecret(path string) (*vault.Secret, error) {
	return f.secret, f.err
}
