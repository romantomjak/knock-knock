package main

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
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
	value := string(pair.Value)
	return value, nil
}
