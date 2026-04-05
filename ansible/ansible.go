// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"encoding/json"
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

// MarshalJSON implements the json.Marshaler interface for Item
func (i Inventory) MarshalJSON() ([]byte, error) {
	// 1. Marshal the struct fields (excluding Metadata due to `json:"-"` tag)
	// Use an anonymous struct or type alias to avoid infinite recursion
	type Alias Inventory
	// Create a copy for marshaling purposes
	aux := struct {
		Alias
	}{
		Alias: Alias{
			Meta: i.Meta,
			All:  i.All,
		},
	}

	// Marshal the auxiliary struct
	msg1, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}

	// 2. Marshal the Metadata map
	msg2, err := json.Marshal(i.Groups)
	if err != nil {
		return nil, err
	}

	// 3. Combine the two JSON objects
	// msg1 is like {"name":"Widget","price":9.99}
	// msg2 is like {"color":"red","size":"large"}
	// We need to combine them into {"name":"Widget","price":9.99,"color":"red","size":"large"}

	// Check if both have content, and combine them
	if len(msg1) > 2 && len(msg2) > 2 { // Check for non-empty objects, e.g., not "{}"
		// Append msg2 (without its opening "{") to msg1 (without its closing "}")
		// and add a comma in between.
		combined := append(msg1[:len(msg1)-1], append([]byte{','}, msg2[1:]...)...)
		return combined, nil
	} else if len(msg2) > 2 {
		return msg2, nil
	}
	return msg1, nil
}
