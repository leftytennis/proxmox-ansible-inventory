// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
)

// GetHosts returns a sorted list of hosts
func (i *Inventory) GetHosts(hosts MapHostVar, excludedHosts mapset.Set[string]) []string {

	// Create a list of hosts
	keys := []string{}

	// Loop through the hosts
	for host := range hosts {
		if !excludedHosts.ContainsOne(host) {
			keys = append(keys, host)
		}
	}

	// Sort the list of hosts
	sort.Strings(keys)
	
	// Return the list of sorted hosts
	return keys
}
