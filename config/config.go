// Package config contains the configuration types for proxmox-ansible-inventory
package config

import "errors"

// CheckRequiredValues checks for required values in the config file
func (p *Params) CheckRequiredValues () error {
	
	if p.Proxmox.API.User == "" {
		return errors.New("proxmox.api.user is required")
	}
	
	if p.Proxmox.API.Token == "" {
		return errors.New("proxmox.api.token is required")
	}
	
	if p.Proxmox.API.Secret == "" {
		return errors.New("proxmox.api.secret is required")
	}
	
	if p.Proxmox.API.URL == "" {
		return errors.New("proxmox.api.url is required")
	}
	
	return nil
}