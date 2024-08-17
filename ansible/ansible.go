// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"sort"
)

// GetHosts returns a sorted list of hosts
func (i *Inventory) GetHosts(hosts MapHostVar, excludedHostsMap map[string]bool) []string {

	keys := []string{}

	for host := range hosts {
		if _, ok := excludedHostsMap[host]; !ok {
			keys = append(keys, host)
		}
	}

	sort.Strings(keys)
	
	return keys
}
