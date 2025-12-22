// Package main is the entry point for the application
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/leftytennis/proxmox-ansible-inventory/ansible"
	"github.com/leftytennis/proxmox-ansible-inventory/config"
	"github.com/leftytennis/proxmox-ansible-inventory/proxmox"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Config is the configuration parameters used by proxmox-ansible-inventory
	Config = config.Params{}
	// ExcludedHosts is a set of excluded hosts
	excludedHosts = mapset.NewSet[string]()
	// GitVersion is the version of the program
	GitVersion = "unknown"
	// GitSha is the git commit hash
	GitSha = "unknown"
	// GitDate is the date the program was built
	GitDate = "unknown"
	// Flags used by this program
	apiToken    string
	baseURL     string
	helpFlag    bool
	listFlag    bool
	versionFlag bool
)

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show program help")
	pflag.BoolVarP(&listFlag, "list", "", true, "list the inventory")
	pflag.BoolVarP(&versionFlag, "version", "", false, "show program version")
}

func main() {

	// Setup viper
	err := setupViper()
	if err != nil {
		fmt.Printf("error setting up viper: %v\n", err)
		os.Exit(1)
	}

	// Parse command line flags
	pflag.Parse()

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

	// Build excluded hosts map
	excludedHosts.Append(Config.Proxmox.Exclude...)

	// Create a new Proxmox client
	ctx := context.Background()

	pm := proxmox.NewClient(&Config)

	// Create proxmox inventory structure
	inv := ansible.Inventory{
		Meta:            ansible.InventoryMeta{},
		All:             ansible.InventoryAll{Children: []string{"containers", "virtual_machines"}},
		Containers:      ansible.InventoryContainers{Hosts: []string{}},
		VirtualMachines: ansible.InventoryVirtualMachines{Hosts: []string{}},
	}

	// Create host vars map
	hostVarMap := make(ansible.MapHostVar)
	inv.Meta.HostVars = hostVarMap

	// Get list of Proxmox nodes
	nodeList, err := pm.GetNodes(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	// Get Proxmox virtual machines and containers from each Proxmox node
	for _, nodeData := range nodeList.Data {

		// Get Proxmox VM list
		vms, err := pm.GetVMList(ctx, nodeData.Node, excludedHosts)
		if err != nil {
			fmt.Printf("error getting proxmox vms: %s\n", err)
			os.Exit(1)
		}
		inv.VirtualMachines.Hosts = append(inv.VirtualMachines.Hosts, vms...)

		// Get Proxmox LXC list
		lxcs, err := pm.GetLxcList(ctx, nodeData.Node, excludedHosts)
		if err != nil {
			fmt.Printf("error getting proxmox lxc list: %s\n", err)
			os.Exit(1)
		}
		inv.Containers.Hosts = append(inv.Containers.Hosts, lxcs...)
	}

	// Pretty print our json inventory
	str, err := json.MarshalIndent(inv, "", "   ")
	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", str)

	os.Exit(0)
}

// setupViper sets up the viper configuration
func setupViper() error {

	// Setup viper
	viper.SetEnvPrefix("PAI")
	viper.SetConfigName(".proxmox-ansible-inventory.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/proxmox-ansible-inventory/")

	// Set defaults
	viper.SetDefault("proxmox.lookup", false)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Unmarshal config values into a Config struct
	err := viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	err = Config.CheckRequiredValues()
	if err != nil {
		return err
	}

	return nil
}
