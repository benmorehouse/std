package hashicorp

import (
	"log"
	"net/http"
	"time"

	"github.com/benmorehouse/std/configs"
	"github.com/benmorehouse/std/repo"
	"github.com/hashicorp/vault/api"
)

type hashicorpClient struct {
	client *api.Client
}

// DefaultVaultClient will return us a client for our secrest stowed away in vault
func DefaultVaultClient() (repo.Repo, error) {
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
func (c *hashicorpClient) Get(key string) string {
	resp, err := c.client.Logical().Read(key)
	if err != nil {
		log.Println(err)
		return ""
	}

	raw, exists := resp.Data[key]
	if !exists {
		log.Println(err)
		return ""
	}

	value, ok := raw.(string)
	if !ok {
		log.Println(err)
		return ""
	}

	return value
}

// Put will go and put a value into vault
func (c *hashicorpClient) Put(key, value string) (err error) {
	_, err = c.client.Logical().Write(key, map[string]interface{}{key: value})
	return
}

// List will go and put a value into vault
func (c *hashicorpClient) List() []string {
	// list, err := c.client.Logical().List("")
	return nil
}

// Remove will go and put a value into vault
func (c *hashicorpClient) Remove(key string) (err error) {
	_, err = c.client.Logical().Delete(key)
	return err
}
