// Package ansible contains the types and methods for implementing the Ansible inventory
package ansible

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
)

// GetHosts returns a sorted list of hosts
func GetHosts(hosts HostVarMap, excludedHosts map[string]bool) []string {
	keys := []string{}
	for host := range hosts {
		if _, ok := excludedHosts[host]; !ok {
			keys = append(keys, host)
		}
	}
	sort.Strings(keys)
	return keys
}

// ReadHosts reads the hosts file and returns a slice of hosts
func ReadHosts(excludedHosts map[string]bool) (HostVarMap, error) {

	hosts := make(HostVarMap)

	f, err := os.Open("/Users/jefft/develop/ansible/hosts")
	defer f.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			continue
		}
		re := regexp.MustCompile(`\S+`)
		args := re.FindAllString(line, -1)
		host := args[0]
		if _, ok := excludedHosts[host]; !ok {
			for i, arg := range args {
				if i == 0 {
					host = arg
					hosts[host] = make(map[string]string)
				} else {
					key := strings.Split(arg, "=")[0]
					value := strings.Split(arg, "=")[1]
					hosts[host][key] = value
				}
			}
		}
	}

	return hosts, nil
}
