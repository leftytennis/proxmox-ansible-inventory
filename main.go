// Package main is the entry point for the application
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/leftytennis/proxmox-ansible-inventory/ansible"
	"github.com/leftytennis/proxmox-ansible-inventory/proxmox"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Config is the configuration parameters used by proxmox-ansible-inventory
	Config = ConfigParams{}
	// ExcludedHostsMap is a map of excluded hosts
	excludedHostsMap = make(map[string]bool)
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

// ConfigParams is the configuration info used by proxmox-ansible-inventory
type ConfigParams struct {
	Proxmox ConfigProxmox `mapstructure:"proxmox"`
}

// ConfigProxmox is the Proxmox section of the config file
type ConfigProxmox struct {
	// APIToken is the Proxmox API token
	APIToken string `mapstructure:"api_token"`
	// BaseURL is the Proxmox API base URL
	BaseURL string `mapstructure:"base_url"`
	// Exclude is a list of hostnames to exclude from the inventory
	Exclude []string `mapstructure:"exclude"`
	// Lookup is a flag to enable or disable IP address lookup
	Lookup bool `mapstructure:"lookup"`
}

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

	// Check for required flags
	// if len(hostFlag) == 0 && !listFlag {
	// 	fmt.Printf("one of --list or --host is required\n")
	// }

	// Build excluded hosts map
	if len(Config.Proxmox.Exclude) > 0 {
		for _, host := range Config.Proxmox.Exclude {
			excludedHostsMap[host] = true
		}
	}

	// Create a new Proxmox client
	ctx := context.Background()

	pm := proxmox.NewClient(Config.Proxmox.BaseURL, Config.Proxmox.APIToken)

	// Create proxmox inventory structure
	inv := ansible.Inventory{
		Meta: ansible.InventoryMeta{},
		All:  ansible.InventoryAll{Children: []string{"containers", "virtual_machines"}},
		Lxcs: ansible.InventoryLxcs{Hosts: []string{}},
		VMs:  ansible.InventoryVMs{Hosts: []string{}},
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
		vms, err := pm.GetVMList(ctx, nodeData.Node, excludedHostsMap)
		if err != nil {
			fmt.Printf("error getting proxmox vms: %s\n", err)
			os.Exit(1)
		}
		inv.VMs.Hosts = append(inv.VMs.Hosts, vms...)
		if Config.Proxmox.Lookup {
			for _, vm := range vms {
				_, err := inv.LookupIPAddress(hostVarMap, excludedHostsMap, vm)
				if err != nil {
					fmt.Printf("error getting ip address for host %s: %s\n", vm, err)
				}
			}
		}

		// Get Proxmox LXC list
		lxcs, err := pm.GetLxcList(ctx, nodeData.Node, excludedHostsMap)
		if err != nil {
			fmt.Printf("error getting proxmox lxc list: %s\n", err)
			os.Exit(1)
		}
		inv.Lxcs.Hosts = append(inv.Lxcs.Hosts, lxcs...)
		if Config.Proxmox.Lookup {
			for _, lxc := range lxcs {
				_, err := inv.LookupIPAddress(hostVarMap, excludedHostsMap, lxc)
				if err != nil {
					fmt.Printf("error getting ip address for host %s: %s\n", lxc, err)
				}
			}
		}
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

	// Proxmox api_token is required
	if !viper.IsSet("proxmox.api_token") {
		return errors.New("api_token is required and is not found in the config file")
	}

	// Proxmox base_url is required
	if !viper.IsSet("proxmox.base_url") {
		return errors.New("base_url is required and is not found in the config file")
	}


	// Unmarshal config values into a Config struct
	err := viper.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}
