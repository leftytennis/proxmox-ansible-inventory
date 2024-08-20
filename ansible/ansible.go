// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"net"
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

// LookupIPAddress returns the IP addresses for a host
func (i *Inventory) LookupIPAddress(mapHostVar MapHostVar, excludedHostsMap map[string]bool, host string) ([]string, error) {

	// Lookup the IP address for the host
	ips, err := net.LookupIP(host)
	if err != nil {
		return []string{}, err
	}

	ipList := []string{}
	
	// Print the IP addresses
	for _, ip := range ips {
		ipList = append(ipList, ip.String())
	}

	if _, ok := excludedHostsMap[host]; !ok {
		mapHostVar[host] = map[string]string{"ansible_host": ipList[0]}
		return ipList, nil
	}
	
	// Return the first IP address
	return []string{}, nil
}
