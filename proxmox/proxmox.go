// Package proxmox provides a client for the Proxmox API.
package proxmox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NewClient creates a new Client
func NewClient(baseURL string, apiToken string) *Client {
	return &Client{
		BaseURL:    baseURL,
		apiToken:   "PVEAPIToken=" + apiToken,
		HTTPClient: &http.Client{Timeout: time.Second * 30},
	}
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", c.apiToken)
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
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseURL, url), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// GetLxcConfig performs a GET request to the Proxmox API
func (c *Client) GetLxcConfig(ctx context.Context, node string, vmid int) (*LxcConfig, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/lxc/%d/config", c.BaseURL, node, vmid), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &LxcConfig{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetLxcList returns a lsit of LXC containers
func (c *Client) GetLxcList(ctx context.Context, node string, excludedHosts map[string]bool) ([]string, error) {

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

	data := &LxcResponse{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	lxcs := []string{}

	for _, lxc := range data.Data {
		if lxc.Status == "running" {
			if _, ok := excludedHosts[lxc.Name]; !ok {
				lxcs = append(lxcs, lxc.Name)
			}
		}
	}

	return lxcs, nil
}

// GetLxcs performs a GET request to the Proxmox API
func (c *Client) GetLxcs(ctx context.Context, node string) (*LxcResponse, error) {

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

	data := &LxcResponse{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetNodes performs a GET request to the Proxmox API
func (c *Client) GetNodes(ctx context.Context) (*NodeList, error) {

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

	data := &NodeList{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetQemuNetworkConfig performs a GET request to the Proxmox API
func (c *Client) GetQemuNetworkConfig(ctx context.Context, node string, vmid int) (*QemuAgentNetworkResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu/%d/agent/network-get-interfaces", c.BaseURL, node, vmid), nil)
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
func (c *Client) GetSubdirs(ctx context.Context) (*Subdir, error) {

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

	data := &Subdir{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVersion performs a GET request to the Proxmox API
func (c *Client) GetVersion(ctx context.Context) (*Version, error) {

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

	data := &Version{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVMConfig performs a GET request to the Proxmox API
func (c *Client) GetVMConfig(ctx context.Context, node string, vmid int) (*VMConfig, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu/%d/config", c.BaseURL, node, vmid), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data := &VMConfig{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetVMList returns a list of VMs
func (c *Client) GetVMList(ctx context.Context, node string, excludedHosts map[string]bool) ([]string, error) {

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

	data := &VMList{}

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
func (c *Client) GetVMs(ctx context.Context, node string) (*VMList, error) {

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

	data := &VMList{}

	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
