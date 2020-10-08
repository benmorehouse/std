package hashicorp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/benmorehouse/std/configs"
	"github.com/hashicorp/vault/api"
)

// VaultClient is the vault client that will provide super
// secret awesome access to stuff
type VaultClient interface {
	Get(key string) (string, error)
	Put(key, value string) error
}

type hashicorpClient struct {
	client *api.Client
}

// DefaultVaultClient will return us a client for our secrest stowed away in vault
func DefaultVaultClient() (VaultClient, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := api.NewClient(&api.Config{
		Address:    configs.STDConf.VaultAddr,
		HttpClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	client.SetToken(configs.STDConf.VaultStaticKey)
	return &hashicorpClient{client: client}, nil
}

// Get will go into our vault instance and get a secret
func (c *hashicorpClient) Get(key string) (string, error) {
	resp, err := c.client.Logical().Read(key)
	if err != nil {
		return "", nil
	}

	raw, exists := resp.Data[key]
	if !exists {
		return "", fmt.Errorf("password_not_part_of_data: %s", err.Error())
	}

	value, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("value_invalid_type: %s", err.Error())
	}

	return value, nil
}

// Put will go and put a value into vault
func (c *hashicorpClient) Put(key, value string) (err error) {
	_, err = c.client.Logical().Write(key, map[string]interface{}{key: value})
	return
}
