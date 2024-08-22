// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"sort"
)

// GetHosts returns a sorted list of hosts
func (i *Inventory) GetHosts(hosts MapHostVar, excludedHostsMap map[string]bool) []string {

	// Create a list of hosts
	keys := []string{}

	// Loop through the hosts
	for host := range hosts {
		if _, ok := excludedHostsMap[host]; !ok {
			keys = append(keys, host)
		}
	}

	// Sort the list of hosts
	sort.Strings(keys)
	
	// Return the list of sorted hosts
	return keys
}
