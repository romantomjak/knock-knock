package main

// Client is an interface for communicating with remote services like Consul or Vault
type Client interface {
	Read(path string) (interface{}, error)
}
