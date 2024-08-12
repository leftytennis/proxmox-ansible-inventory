// Package main is the entry point for the application
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/leftytennis/proxmox-ansible-inventory/ansible"
	"github.com/leftytennis/proxmox-ansible-inventory/proxmox"
)

const (
	proxmoxAPIKey = "PVEAPIToken=root@pam!pmnbsync=afa1e088-8ed8-45e4-8b33-2d0b6f456959"
)

var (
	// ExcludeHosts is a list of hosts to exclude from the inventory
	excludeHosts = map[string]bool{"pfsense-lab": true, "haos1": true, "mba01": true}
	version      = "unknown"
	commit       = "unknownx"
	date         = "unknown"
	versionFlag  bool
	helpFlag     bool
)

func init() {
	flag.BoolVar(&helpFlag, "help", false, "show program help")
	flag.BoolVar(&versionFlag, "version", false, "show program version")
}

func main() {

	flag.Parse()

	if versionFlag {
		lenCommit := len(commit)
		min := int(math.Min(7, float64(lenCommit)))
		fmt.Printf("%s (v%s git %s built on %s)\n", os.Args[0], version, commit[:min], date)
		os.Exit(0)
	}

	if helpFlag {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	ctx := context.Background()

	pm := proxmox.NewProxmoxClient(proxmoxAPIKey)

	inv := ansible.Inventory{
		Meta: ansible.InventoryMeta{},
		All:  ansible.InventoryAll{Children: []string{"proxmox_lxcs", "proxmox_vms", "static"}},
		LXCs: ansible.InventoryProxmoxLXCs{},
		VMs:  ansible.InventoryProxmoxVMs{},
	}

	hosts, err := ansible.ReadHosts(excludeHosts)

	if err != nil {
		fmt.Printf("error reading hosts file: %v\n", err)
		os.Exit(1)
	}

	inv.Static.Hosts = ansible.GetHosts(hosts, excludeHosts)
	inv.Meta.HostVars = hosts

	// // Get Proxmox nodes
	// nodes, err := pm.GetNodes(ctx)

	// if err != nil {
	// 	fmt.Printf("error: %s\n", err)
	// 	os.Exit(1)
	// }

	// Get Proxmox VM list

	vms, err := pm.GetVMList(ctx, "pve1", excludeHosts)

	if err != nil {
		fmt.Printf("error getting proxmox vms: %s\n", err)
		os.Exit(1)
	}

	inv.VMs.Hosts = append(inv.VMs.Hosts, vms...)

	// Get Proxmox LXC list

	lxcs, err := pm.GetLXCList(ctx, "pve1")

	if err != nil {
		fmt.Printf("error getting proxmox lxc list: %s\n", err)
		os.Exit(1)
	}

	inv.LXCs.Hosts = append(inv.LXCs.Hosts, lxcs...)

	// Pretty print our json inventory

	str, err := json.MarshalIndent(inv, "", "    ")

	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", str)

}
