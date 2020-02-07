package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/api"
	vaultapi "github.com/hashicorp/vault/api"
)

// Client is an interface for communicating with remote services like Consul or Vault
type Client interface {
	Read(path string) (interface{}, error)
}

// Consul is a wrapper around a real Consul API client
type Consul struct {
	client *consulapi.Client
}

// NewConsulClient creates a new Consul API client
func NewConsulClient() (*Consul, error) {
	consulConfig := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("client: consul: %s", err)
	}
	return &Consul{client}, nil
}

// Read queries the Consul API
func (c *Consul) Read(path string) (interface{}, error) {
	pair, _, err := c.client.KV().Get(path, nil)
	if err != nil {
		return nil, fmt.Errorf("client: consul: %s", err)
	}
	if pair == nil {
		return nil, fmt.Errorf("client: consul: key %q does not exist", path)
	}
	value := string(pair.Value)
	return value, nil
}

// Vault is a wrapper around a real Consul API client
type Vault struct {
	client *vaultapi.Client
}

// NewVaultClient creates a new Vault API client
func NewVaultClient() (*Vault, error) {
	vaultConfig := vaultapi.DefaultConfig()
	client, err := vaultapi.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("client: vault: %s", err)
	}

	// try vault token
	method := ""
	token := os.Getenv("VAULT_TOKEN")

	// next try github token
	if token == "" {
		method = "github"
		token = os.Getenv("VAULT_AUTH_GITHUB_TOKEN")
	}

	// bail if token is still empty
	if token == "" {
		return nil, fmt.Errorf("client: vault: no auth token provided")
	}

	// call to get a token
	if method == "github" {
		secret, err := client.Logical().Write("auth/github/login", map[string]interface{}{
			"token": strings.TrimSpace(token),
		})
		if err != nil {
			return nil, fmt.Errorf("client: vault: %s", err)
		}
		if secret == nil {
			return nil, fmt.Errorf("client: vault: empty response from credential provider")
		}
		if secret.Auth == nil {
			return nil, fmt.Errorf("client: vault: no secret auth")
		}
		if secret.Auth.ClientToken == "" {
			return nil, fmt.Errorf("client: vault: no token returned")
		}
		token = secret.Auth.ClientToken
	}

	client.SetToken(token)

	return &Vault{client}, nil
}

// Read queries the Vault API
func (c *Vault) Read(p string) (interface{}, error) {
	secretPath := p

	mountPath, isKVv2, err := c.isKVv2(p)
	if err != nil {
		isKVv2 = false
	} else if isKVv2 {
		secretPath = c.addPrefixToKVPath(secretPath, mountPath, "data")
	}

	secret, err := c.client.Logical().ReadWithData(secretPath, nil)
	if err != nil {
		return nil, fmt.Errorf("client: vault: %s", err)
	}
	return secret, nil
}

func (c *Vault) isKVv2(path string) (string, bool, error) {
	r := c.client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := c.client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// if we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", false, nil
		}

		// anonymous requests may fail to access /sys/internal/ui path
		// Vault v1.1.3 returns 500 status code but may return 4XX in future
		if c.client.Token() == "" {
			return "", false, nil
		}

		return "", false, err
	}
	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", false, err
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	var mountType string
	if mountTypeRaw, ok := secret.Data["type"]; ok {
		mountType = mountTypeRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, false, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, false, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, false, nil
	case "2":
		return mountPath, mountType == "kv", nil
	}

	return mountPath, false, nil
}

func (c *Vault) addPrefixToKVPath(p, mountPath, apiPrefix string) string {
	switch {
	case p == mountPath, p == strings.TrimSuffix(mountPath, "/"):
		return path.Join(mountPath, apiPrefix)
	default:
		p = strings.TrimPrefix(p, mountPath)
		// don't add /data to the path if it's been added manually.
		if strings.HasPrefix(p, apiPrefix) {
			return path.Join(mountPath, p)
		}
		return path.Join(mountPath, apiPrefix, p)
	}
}
