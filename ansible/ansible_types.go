// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

// Inventory is the top-level structure for the Ansible inventory
type Inventory struct {
	Meta   InventoryMeta   `json:"_meta"`
	All    InventoryAll    `json:"all"`
	Lxcs   InventoryLxcs   `json:"lxc_containers"`
	VMs    InventoryVMs    `json:"virtual_machines"`
	Static InventoryStatic `json:"static"`
}

// InventoryMeta is the metadata for the Ansible inventory
type InventoryMeta struct {
	HostVars MapHostVar `json:"hostvars"`
}

// MapHostVar is a map of ansible host variables
type MapHostVar map[string]map[string]string

// InventoryAll is the "all" group in the Ansible inventory
type InventoryAll struct {
	Children []string `json:"children"`
}

// InventoryLxcs is the "proxmox_lxcs" group in the Ansible inventory
type InventoryLxcs struct {
	Hosts []string `json:"hosts"`
}

// InventoryVMs is the "proxmox_vms" group in the Ansible inventory
type InventoryVMs struct {
	Hosts []string `json:"hosts"`
}

// InventoryStatic is the "static" group (not proxmox lxc or vm) in the Ansible inventory
type InventoryStatic struct {
	Hosts []string `json:"hosts"`
}
