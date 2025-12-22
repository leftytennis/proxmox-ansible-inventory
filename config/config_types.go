// Package config contains the configuration types for proxmox-ansible-inventory
package config

// Params is the configuration info used by proxmox-ansible-inventory
type Params struct {
	Proxmox ProxmoxParams `mapstructure:"proxmox"`
}

// ProxmoxParams is the Proxmox section of the config file
type ProxmoxParams struct {
	// APIParams is the Proxmox API token
	API APIParams `mapstructure:"api"`
	// Exclude is a list of hostnames to exclude from the inventory
	Exclude []string `mapstructure:"exclude"`
}

// APIParams is the api_token section of the config file
type APIParams struct {
	// Secret is the api token secret
	Secret string `mapstructure:"secret"`
	// Token is the api token
	Token  string `mapstructure:"token"`
	// URL is the Proxmox API base URL
	URL string `mapstructure:"url"`
	// User is the api token user
	User   string `mapstructure:"user"`
}
