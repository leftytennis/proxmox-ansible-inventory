// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

// Inventory is the top-level structure for the Ansible inventory
type Inventory struct {
	Meta   InventoryMeta     `json:"_meta"`
	All    InventoryAll      `json:"all"`
	Groups InventoryGroupMap `json:"-"`
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

// InventoryGroup is a single Ansible inventory group
type InventoryGroup struct {
	Hosts []string `json:"hosts"`
}

// InventoryGroupMap is a map of inventory groups to their hosts
type InventoryGroupMap map[string]InventoryGroup
