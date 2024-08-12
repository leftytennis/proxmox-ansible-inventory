// Package main is the entry point for the application
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	commit       = "unknown"
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
		fmt.Printf("%s (v%s built on %s)\n", os.Args[0], version, date)
		os.Exit(0)
	}

	if helpFlag {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	ctx := context.Background()

	pm := proxmox.NewClient(proxmoxAPIKey)

	inv := ansible.Inventory{
		Meta: ansible.InventoryMeta{},
		All:  ansible.InventoryAll{Children: []string{"lxc_containers", "virtual_machines", "static"}},
		Lxcs: ansible.InventoryLxcs{},
		VMs:  ansible.InventoryVMs{},
	}

	hostVarMap, err := inv.ReadHosts(excludeHosts)

	if err != nil {
		fmt.Printf("error reading hosts file: %v\n", err)
		os.Exit(1)
	}

	inv.Static.Hosts = inv.GetHosts(hostVarMap, excludeHosts)
	inv.Meta.HostVars = hostVarMap

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

	lxcs, err := pm.GetLxcList(ctx, "pve1")

	if err != nil {
		fmt.Printf("error getting proxmox lxc list: %s\n", err)
		os.Exit(1)
	}

	inv.Lxcs.Hosts = append(inv.Lxcs.Hosts, lxcs...)

	// Pretty print our json inventory

	str, err := json.MarshalIndent(inv, "", "    ")

	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", str)

}
