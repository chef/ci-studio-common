package account

import (
	"fmt"
	"reflect"

	"github.com/chef/ci-studio-common/lib"
	"github.com/chef/ci-studio-common/vault-util/internal/vault"
	"github.com/spf13/viper"
)

type (
	// AWSAccount holds credentials for AWS account
	AWSAccount struct {
		Name            string
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
		Region          string
	}
)

var (
	// AWSUserFieldMap maps snake case fields to camelcase struct keys
	AWSUserFieldMap = map[string]string{
		"access_key_id":     "AccessKeyID",
		"secret_access_key": "SecretAccessKey",
		"session_token":     "SessionToken",
		"region":            "Region",
	}
)

// InitAWSConfig is called through the chain to establish global defaults
func InitAWSConfig() {
	viper.SetDefault("aws.region", "us-east-1")
	viper.SetDefault("aws.dynamic_mount_root", "account/dynamic/aws")
}

// NewAWSAccount returns new AWSAccount struct
func NewAWSAccount(name string) (*AWSAccount, error) {
	vaultClient, err := vault.NewClient()
	if err != nil {
		return nil, err
	}

	secret, err := vaultClient.Read(fmt.Sprintf("%s/%s/sts/default", viper.GetString("aws.dynamic_mount_root"), name))
	if err != nil {
		return nil, err
	}

	return &AWSAccount{
		Name:            name,
		AccessKeyID:     secret.Data["access_key"].(string),
		SecretAccessKey: secret.Data["secret_key"].(string),
		SessionToken:    secret.Data["security_token"].(string),
		Region:          viper.GetString("aws.region"),
	}, nil
}

func (c *Cache) fetchOrInitAWSAccount(name string) (*AWSAccount, error) {
	var acct *AWSAccount
	var cached bool
	var err error

	if acct, cached = c.AWSAccounts[name]; !cached {
		acct, err = NewAWSAccount(name)

		if err != nil {
			return nil, err
		}

		c.AWSAccounts[name] = acct
	}

	return acct, nil
}

// ConfigureAWS runs commands to configure the given account
func (c *Cache) ConfigureAWS(name string) error {
	acct, err := c.fetchOrInitAWSAccount(name)
	if err != nil {
		return err
	}

	err = lib.ShellOut("aws", "configure", "set", "aws_access_key_id", acct.AccessKeyID, "--profile", acct.Name).Run()
	if err != nil {
		return err
	}

	err = lib.ShellOut("aws", "configure", "set", "aws_secret_access_key", acct.SecretAccessKey, "--profile", acct.Name).Run()
	if err != nil {
		return err
	}

	err = lib.ShellOut("aws", "configure", "set", "aws_session_token", acct.SessionToken, "--profile", acct.Name).Run()
	if err != nil {
		return err
	}

	return lib.ShellOut("aws", "configure", "set", "region", acct.Region, "--profile", acct.Name).Run()
}

// GetAWSSecret provides a mechanism to quickly fetch field
func (c *Cache) GetAWSSecret(name string, field string) (string, error) {
	acct, err := c.fetchOrInitAWSAccount(name)
	if err != nil {
		return "", err
	}

	r := reflect.ValueOf(acct)
	camelField := AWSUserFieldMap[field]
	value := reflect.Indirect(r).FieldByName(camelField)

	return value.String(), nil
}
