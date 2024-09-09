// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

// Inventory is the top-level structure for the Ansible inventory
type Inventory struct {
	Meta            InventoryMeta            `json:"_meta"`
	All             InventoryAll             `json:"all"`
	Containers      InventoryContainers      `json:"containers"`
	VirtualMachines InventoryVirtualMachines `json:"virtual_machines"`
	// Ungrouped InventoryUngrouped `json:"ungrouped"`
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

// InventoryContainers is the "containers" group in the Ansible inventory
type InventoryContainers struct {
	Hosts []string `json:"hosts"`
}

// InventoryVirtualMachines is the "virtual_machines" group in the Ansible inventory
type InventoryVirtualMachines struct {
	Hosts []string `json:"hosts"`
}
