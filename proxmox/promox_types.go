// Package proxmox is a package to interact with the Proxmox VE API
package proxmox

import "net/http"

// LxcConfig is the response for the Proxmox API LXC config
type LxcConfig struct {
	Data LxcConfigData `json:"data"`
}

// LxcConfigData is the struct for the Proxmox API LXC config data
type LxcConfigData struct {
	Swap         int    `json:"swap"`
	Unprivileged int    `json:"unprivileged"`
	Net0         string `json:"net0"`
	Net1         string `json:"net1"`
	Net2         string `json:"net2"`
	Net3         string `json:"net3"`
	Net4         string `json:"net4"`
	Memory       int    `json:"memory"`
	Digest       string `json:"digest"`
	Features     string `json:"features"`
	Description  string `json:"description"`
	Tags         string `json:"tags"`
	Ostype       string `json:"ostype"`
	Rootfs       string `json:"rootfs"`
	Cores        int    `json:"cores"`
	Onboot       int    `json:"onboot"`
	Hostname     string `json:"hostname"`
	Arch         string `json:"arch"`
}

// LxcResponse is the list of Proxmox LXC containers
type LxcResponse struct {
	Data []LxcData `json:"data"`
}

// LxcData is the struct for a Proxmoc LXC container
type LxcData struct {
	Name      string  `json:"name"`
	Disk      int     `json:"disk"`
	Maxdisk   int64   `json:"maxdisk"`
	Type      string  `json:"type"`
	Maxmem    int     `json:"maxmem"`
	Status    string  `json:"status"`
	Pid       int     `json:"pid"`
	Vmid      int     `json:"vmid"`
	Netout    int     `json:"netout"`
	Netin     int     `json:"netin"`
	Diskread  int     `json:"diskread"`
	Uptime    int     `json:"uptime"`
	Diskwrite int     `json:"diskwrite"`
	Tags      string  `json:"tags,omitempty"`
	Maxswap   int     `json:"maxswap"`
	Swap      int     `json:"swap"`
	Mem       int     `json:"mem"`
	Cpus      int     `json:"cpus"`
	CPU       float64 `json:"cpu"`
}

// Client is the struct for the Proxmox API client
type Client struct {
	BaseURL    string
	apiToken   string
	HTTPClient *http.Client
}

// NodeList is the struct for the Proxmox API data
type NodeList struct {
	Data []NodeData `json:"data"`
}

// NodeData is the struct for the Proxmox API node data
type NodeData struct {
	Node           string  `json:"node"`
	Mem            int64   `json:"mem"`
	Level          string  `json:"level"`
	Status         string  `json:"status"`
	Disk           int64   `json:"disk"`
	CPU            float64 `json:"cpu"`
	ID             string  `json:"id"`
	Maxdisk        int64   `json:"maxdisk"`
	SslFingerprint string  `json:"ssl_fingerprint"`
	Uptime         int     `json:"uptime"`
	Maxcpu         int     `json:"maxcpu"`
	Maxmem         int64   `json:"maxmem"`
	Type           string  `json:"type"`
}

// Subdir is the struct for the Proxmox API data
type Subdir struct {
	Data []SubdirData `json:"data"`
}

// SubdirData is the struct for the Proxmox API subdir
type SubdirData struct {
	Subdir string `json:"subdir"`
}

// Version is the struct for the Proxmox API data
type Version struct {
	Data VersionData `json:"data"`
}

// VersionData is the struct for the Proxmox API version info
type VersionData struct {
	Release string `json:"release"`
	Version string `json:"version"`
	RepoID  string `json:"repoid"`
}

// VMConfig is the struct for the Proxmox API VM config
type VMConfig struct {
	Data VMConfigData `json:"data"`
}

// VMConfigData is the struct for the Proxmox API VM config data
type VMConfigData struct {
	Scsihw  string `json:"scsihw"`
	Ide2    string `json:"ide2"`
	Cores   int    `json:"cores"`
	Vmgenid string `json:"vmgenid"`
	CPU     string `json:"cpu"`
	Meta    string `json:"meta"`
	Scsi1   string `json:"scsi1"`
	Agent   string `json:"agent"`
	Digest  string `json:"digest"`
	Numa    int    `json:"numa"`
	Memory  string `json:"memory"`
	Boot    string `json:"boot"`
	Net0    string `json:"net0"`
	Net1    string `json:"net1"`
	Net2    string `json:"net2"`
	Net3    string `json:"net3"`
	Net4    string `json:"net4"`
	Ostype  string `json:"ostype"`
	Name    string `json:"name"`
	Tags    string `json:"tags"`
	Onboot  int    `json:"onboot"`
	Smbios1 string `json:"smbios1"`
	Sockets int    `json:"sockets"`
	Scsi0   string `json:"scsi0"`
}

// VMList is the struct for the Proxmox API data
type VMList struct {
	Data []VM `json:"data"`
}

// VM is the struct for a proxmoc Qemu VM
type VM struct {
	Disk      int     `json:"disk"`
	Name      string  `json:"name"`
	Maxmem    int64   `json:"maxmem"`
	Pid       int     `json:"pid"`
	Status    string  `json:"status"`
	Maxdisk   int64   `json:"maxdisk"`
	Netout    int     `json:"netout"`
	Vmid      int     `json:"vmid"`
	Diskread  int     `json:"diskread"`
	Uptime    int     `json:"uptime"`
	Diskwrite int     `json:"diskwrite"`
	Netin     int     `json:"netin"`
	Tags      string  `json:"tags"`
	Cpus      int     `json:"cpus"`
	CPU       float64 `json:"cpu"`
	Mem       int     `json:"mem"`
}

// QemuAgentIPAddresses is the struct for Qemu IP addresses
type QemuAgentIPAddresses struct {
	IPAddress     string `json:"ip-address"`
	Prefix        int    `json:"prefix"`
	IPAddressType string `json:"ip-address-type"`
}

// QemuAgentNetworkResult is the struct for the Proxmox API data
type QemuAgentNetworkResult struct {
	Name            string                     `json:"name"`
	Statistics      QemuAgentNetworkStatistics `json:"statistics"`
	IPAddresses     []QemuAgentIPAddresses     `json:"ip-addresses"`
	HardwareAddress string                     `json:"hardware-address"`
}

// QemuAgentNetworkData is the struct for the Proxmox API data
type QemuAgentNetworkData struct {
	Result []QemuAgentNetworkResult `json:"result"`
}

// QemuAgentNetworkResponse is the struct for Qemu API response:
// /api2/json/nodes/pve1/qemu/100/agent/network-get-interfaces
type QemuAgentNetworkResponse struct {
	Data QemuAgentNetworkData `json:"data"`
}

// QemuAgentNetworkStatistics is the struct for the Proxmox API network statistics
type QemuAgentNetworkStatistics struct {
	RxBytes   int `json:"rx-bytes"`
	TxDropped int `json:"tx-dropped"`
	TxBytes   int `json:"tx-bytes"`
	RxDropped int `json:"rx-dropped"`
	TxPackets int `json:"tx-packets"`
	TxErrs    int `json:"tx-errs"`
	RxPackets int `json:"rx-packets"`
	RxErrs    int `json:"rx-errs"`
}
