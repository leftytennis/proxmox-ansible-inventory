// Package proxmox is a package to interact with the Proxmox VE API
package proxmox

import "net/http"

// LXCConfig is the response for the Proxmox API LXC config
type LXCConfig struct {
	Data LXCConfigData `json:"data"`
}

// LXCConfigData is the struct for the Proxmox API LXC config data
type LXCConfigData struct {
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

// LXCResponse is the list of Proxmox LXC containers
type LXCResponse struct {
	Data []LXCData `json:"data"`
}

// LXCData is the struct for a Proxmoc LXC container
type LXCData struct {
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

// ProxmoxClient is the struct for the Proxmox API client
type ProxmoxClient struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

// ProxmoxNodeList is the struct for the Proxmox API data
type ProxmoxNodeList struct {
	Data []ProxmoxNodeData `json:"data"`
}

// ProxmoxNodeData is the struct for the Proxmox API node data
type ProxmoxNodeData struct {
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

// ProxmoxSubdir is the struct for the Proxmox API data
type ProxmoxSubdir struct {
	Data []ProxmoxSubdirData `json:"data"`
}

// ProxmoxSubdirData is the struct for the Proxmox API subdir
type ProxmoxSubdirData struct {
	Subdir string `json:"subdir"`
}

// ProxmoxVersion is the struct for the Proxmox API data
type ProxmoxVersion struct {
	Data ProxmoxVersionData `json:"data"`
}

// ProxmoxVersionData is the struct for the Proxmox API version info
type ProxmoxVersionData struct {
	Release string `json:"release"`
	Version string `json:"version"`
	RepoID  string `json:"repoid"`
}

// ProxmoxVMConfig is the struct for the Proxmox API VM config
type ProxmoxVMConfig struct {
	Data ProxmoxVMConfigData `json:"data"`
}

// ProxmoxVMConfigData is the struct for the Proxmox API VM config data
type ProxmoxVMConfigData struct {
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

// ProxmoxVMs is the struct for the Proxmox API data
type ProxmoxVMs struct {
	Data []ProxmoxVM `json:"data"`
}

// ProxmoxVM is the struct for a proxmoc Qemu VM
type ProxmoxVM struct {
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
