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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// ExcludedHostsMap is a map of excluded hosts
	excludedHostsMap = make(map[string]bool)
	// GitVersion is the version of the program
	GitVersion = "unknown"
	// GitSha is the git commit hash
	GitSha = "unknown"
	// GitDate is the date the program was built
	GitDate = "unknown"
	// Flags used by this program
	apiToken      string
	baseURL       string
	excludedHosts []string
	helpFlag      bool
	hostFlag      string
	listFlag      bool
	versionFlag   bool
)

// Config is the configuration info used by proxmox-ansible-inventory
type Config struct {
	// APIToken is the Proxmox API token
	APIToken string
	// BaseURL is the Proxmox API base URL
	BaseURL  string
	// Exclude is a list of hostnames to exclude from the inventory
	Exclude  []string
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.StringSliceVarP(&excludedHosts, "exclude", "e", []string{}, "exclude the specified `hostnames` from the inventory")
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show program help")
	pflag.StringVarP(&hostFlag, "host", "", "", "show info for the specified `hostname`")
	pflag.BoolVarP(&listFlag, "list", "", true, "list the inventory")
	pflag.BoolVarP(&versionFlag, "version", "", false, "show program version")
	pflag.CommandLine.MarkHidden("exclude")
}

func main() {

	// Setup viper
	viper.SetEnvPrefix("PAI")
	viper.SetConfigName(".proxmox-ansible-inventory")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// viper.AutomaticEnv()
	// viper.BindEnv("api_token")
	// viper.BindEnv("base_url")
	// viper.BindEnv("exclude")

	err := viper.ReadInConfig()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Parse command line flags
	pflag.Parse()

	// Get config values from viper
	apiToken := viper.GetString("proxmox.api_token")
	baseURL := viper.GetString("proxmox.base_url")
	excludedHosts := viper.GetStringSlice("proxmox.exclude")

	// Show help if requested
	if helpFlag {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		pflag.PrintDefaults()
		os.Exit(0)
	}

	// Show version if requested
	if versionFlag {
		fmt.Printf("%s (%s built on %s)\n", os.Args[0], GitVersion, GitDate)
		os.Exit(0)
	}

	if len(hostFlag) == 0 && !listFlag {
		fmt.Printf("one of --list or --host is required\n")
	}

	if len(excludedHosts) > 0 {
		for _, host := range excludedHosts {
			excludedHostsMap[host] = true
		}
	}

	ctx := context.Background()

	pm := proxmox.NewClient(baseURL, apiToken)

	inv := ansible.Inventory{
		Meta: ansible.InventoryMeta{},
		All:  ansible.InventoryAll{Children: []string{"containers", "virtual_machines"}},
		Lxcs: ansible.InventoryLxcs{Hosts: []string{}},
		VMs:  ansible.InventoryVMs{Hosts: []string{}},
	}

	hostVarMap := make(ansible.MapHostVar)

	inv.Meta.HostVars = hostVarMap

	// Get list of Proxmox nodes
	nodeList, err := pm.GetNodes(ctx)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	// Get Proxmox VMs and LXCs info from each Proxmox node
	for _, nodeData := range nodeList.Data {

		// Get Proxmox VM list
		vms, err := pm.GetVMList(ctx, nodeData.Node, excludedHostsMap)

		if err != nil {
			fmt.Printf("error getting proxmox vms: %s\n", err)
			os.Exit(1)
		}

		inv.VMs.Hosts = append(inv.VMs.Hosts, vms...)

		// Get Proxmox LXC list
		lxcs, err := pm.GetLxcList(ctx, nodeData.Node, excludedHostsMap)

		if err != nil {
			fmt.Printf("error getting proxmox lxc list: %s\n", err)
			os.Exit(1)
		}

		inv.Lxcs.Hosts = append(inv.Lxcs.Hosts, lxcs...)
	}

	// Pretty print our json inventory

	str, err := json.MarshalIndent(inv, "", "    ")

	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", str)

}
