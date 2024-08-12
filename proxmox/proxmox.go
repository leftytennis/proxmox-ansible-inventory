// Package proxmox provides a client for the Proxmox API.
package proxmox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	proxmoxAPIURL = "https://pve1.lefty.tennis:8006/api2/json/"
)

// NewProxmoxClient creates a new ProxmoxClient
func NewProxmoxClient(apiKey string) *ProxmoxClient {
	return &ProxmoxClient{
		BaseURL:    proxmoxAPIURL,
		apiKey:     apiKey,
		HTTPClient: &http.Client{Timeout: time.Second * 30},
	}
}

func (c *ProxmoxClient) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

// Get performs a GET request to the Proxmox API
func (c *ProxmoxClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseURL, url), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// GetLXCConfig performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetLXCConfig(ctx context.Context, vmid int) (*LXCConfig, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/pve1/lxc/%d/config", c.BaseURL, vmid), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &LXCConfig{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetLXCList returns a lsit of LXC containers
func (c *ProxmoxClient) GetLXCList(ctx context.Context, node string) ([]string, error) {
	
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/lxc", c.BaseURL, node), nil)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &LXCResponse{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	lxcs := []string{}

	for _, lxc := range data.Data {
		if lxc.Status == "running" {
			lxcs = append(lxcs, lxc.Name)
		}
	}

	return lxcs, nil
}

// GetLXCs performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetLXCs(ctx context.Context) (*LXCResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/pve1/lxc", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &LXCResponse{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetNodes performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetNodes(ctx context.Context) (*ProxmoxNodeList, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxNodeList{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetQemuNetworkConfig performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetQemuNetworkConfig(ctx context.Context, vmid int) (*QemuAgentNetworkResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/pve1/qemu/%d/agent/network-get-interfaces", c.BaseURL, vmid), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &QemuAgentNetworkResponse{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetSubdirs performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetSubdirs(ctx context.Context) (*ProxmoxSubdir, error) {

	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxSubdir{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVersion performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetVersion(ctx context.Context) (*ProxmoxVersion, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/version", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxVersion{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVMConfig performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetVMConfig(ctx context.Context, vmid int) (*ProxmoxVMConfig, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/pve1/qemu/%d/config", c.BaseURL, vmid), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxVMConfig{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVMList returns a list of VMs
func (c *ProxmoxClient) GetVMList(ctx context.Context, node string, excludedHosts map[string]bool) ([]string, error) {
	
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu", c.BaseURL, node), nil)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxVMs{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	vms := []string{}

	for _, vm := range data.Data {
		if vm.Status == "running" {
			if _, ok := excludedHosts[vm.Name]; !ok {
				vms = append(vms, vm.Name)
			}
		}
	}

	return vms, nil
}

// GetVMs performs a GET request to the Proxmox API
func (c *ProxmoxClient) GetVMs(ctx context.Context, node string) (*ProxmoxVMs, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu", c.BaseURL, node), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &ProxmoxVMs{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
