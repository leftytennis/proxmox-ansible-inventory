// Package proxmox provides a client for the Proxmox API.
package proxmox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/leftytennis/proxmox-ansible-inventory/config"
)

// NewClient creates a new Client
func NewClient(cfg *config.Params) *Client {

	var apiToken string
	var baseURL string

	apiToken = "PVEAPIToken=" + cfg.Proxmox.API.User + "!" + cfg.Proxmox.API.Token + "=" + cfg.Proxmox.API.Secret
	baseURL = strings.TrimSuffix(cfg.Proxmox.API.URL, "/") + "/api2/json"

	return &Client{
		apiToken:   apiToken,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: time.Second * 30},
	}
}

// doRequest performs the HTTP request
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {

	// Set the required headers
	req.Header.Add("Authorization", c.apiToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Do the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check the status code
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// return the response and no error
	return resp, nil
}

// Get performs a GET request to the Proxmox API
func (c *Client) Get(url string) (*http.Response, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseURL, url), nil)
	if err != nil {
		return nil, err
	}

	// do the request and return the response
	return c.doRequest(req)
}

// GetLxcConfig performs a GET request to the Proxmox API
func (c *Client) GetLxcConfig(ctx context.Context, node string, vmid int) (*LxcConfig, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/lxc/%d/config", c.BaseURL, node, vmid), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the LxcConfig struct
	data := &LxcConfig{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetLxcList returns a lsit of LXC containers
func (c *Client) GetLxcList(ctx context.Context, node string, excludedHosts mapset.Set[string]) ([]string, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/lxc", c.BaseURL, node), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the LxcResponse struct
	data := &LxcResponse{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Create a list of LXC containers
	lxcs := []string{}

	// Loop through the data and append the LXC containers to the list
	for _, lxc := range data.Data {
		if lxc.Status == "running" {
			if !excludedHosts.ContainsOne(lxc.Name) {
				lxcs = append(lxcs, lxc.Name)
			}
		}
	}

	// Return the list of LXC containers and no error
	return lxcs, nil
}

// GetLxcs performs a GET request to the Proxmox API
func (c *Client) GetLxcs(ctx context.Context, node string) (*LxcResponse, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/lxc", c.BaseURL, node), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the LxcResponse struct
	data := &LxcResponse{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetNodes performs a GET request to the Proxmox API
func (c *Client) GetNodes(ctx context.Context) (*NodeList, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the NodeList struct
	data := &NodeList{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetQemuNetworkConfig performs a GET request to the Proxmox API
func (c *Client) GetQemuNetworkConfig(ctx context.Context, node string, vmid int) (*QemuAgentNetworkResponse, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu/%d/agent/network-get-interfaces", c.BaseURL, node, vmid), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the QemuAgentNetworkResponse struct
	data := &QemuAgentNetworkResponse{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetSubdirs performs a GET request to the Proxmox API
func (c *Client) GetSubdirs(ctx context.Context) (*Subdir, error) {

	// Create the request
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the Subdir struct
	data := &Subdir{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetVersion performs a GET request to the Proxmox API
func (c *Client) GetVersion(ctx context.Context) (*Version, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/version", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the Version struct
	data := &Version{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetVMConfig performs a GET request to the Proxmox API
func (c *Client) GetVMConfig(ctx context.Context, node string, vmid int) (*VMConfig, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu/%d/config", c.BaseURL, node, vmid), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the VMConfig struct
	data := &VMConfig{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}

// GetVMList returns a list of VMs
func (c *Client) GetVMList(ctx context.Context, node string, excludedHosts mapset.Set[string]) ([]string, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu", c.BaseURL, node), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the VMList struct
	data := &VMList{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Create a list of VMs
	vms := []string{}

	// Loop through the data and append the VMs to the list
	for _, vm := range data.Data {
		if vm.Status == "running" {
			if !excludedHosts.ContainsOne(vm.Name) {
				vms = append(vms, vm.Name)
			}
		}
	}

	// Return the list of VMs and no error
	return vms, nil
}

// GetVMs performs a GET request to the Proxmox API
func (c *Client) GetVMs(ctx context.Context, node string) (*VMList, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/nodes/%s/qemu", c.BaseURL, node), nil)
	if err != nil {
		return nil, err
	}

	// Add the context
	req = req.WithContext(ctx)

	// Do the request
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	// Close the response body
	defer resp.Body.Close()

	// Create the VMList struct
	data := &VMList{}

	// Decode the response
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	// Return the data and no error
	return data, nil
}
