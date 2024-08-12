// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

// Inventory is the top-level structure for the Ansible inventory
type Inventory struct {
	Meta   InventoryMeta        `json:"_meta"`
	All    InventoryAll         `json:"all"`
	LXCs   InventoryProxmoxLXCs `json:"proxmox_lxcs"`
	VMs    InventoryProxmoxVMs  `json:"proxmox_vms"`
	Static InventoryStatic      `json:"static"`
}

// InventoryMeta is the metadata for the Ansible inventory
type InventoryMeta struct {
	HostVars HostVarMap `json:"hostvars"`
}

// HostVarMap is a map of ansible host variables
type HostVarMap map[string]map[string]string

// InventoryAll is the "all" group in the Ansible inventory
type InventoryAll struct {
	Children []string `json:"children"`
}

// InventoryProxmoxLXCs is the "proxmox_lxcs" group in the Ansible inventory
type InventoryProxmoxLXCs struct {
	Hosts []string `json:"hosts"`
}

// InventoryProxmoxVMs is the "proxmox_vms" group in the Ansible inventory
type InventoryProxmoxVMs struct {
	Hosts []string `json:"hosts"`
}

// InventoryUngrouped is the "ungrouped" group in the Ansible inventory
type InventoryStatic struct {
	Hosts []string `json:"hosts"`
}
