// Package main is the entry point for the application
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"

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
	helpFlag    bool
	hostFlag    string
	listFlag    bool
	versionFlag bool
)

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show program help")
	pflag.StringVarP(&hostFlag, "host", "", "", "show variables for a single host")
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
		Meta: ansible.InventoryMeta{},
		All:  ansible.InventoryAll{Children: []string{"proxmox_lxcs", "proxmox_vms", "ungrouped"}},
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

	type hostInfo struct {
		node string
		vmid int
	}

	roles := make(map[string][]string)
	lxcNames := []string{}
	vmNames := []string{}
	lxcHosts := make(map[string]hostInfo)
	vmHosts := make(map[string]hostInfo)

	// Get Proxmox virtual machines and containers from each Proxmox node
	for _, nodeData := range nodeList.Data {

		// Get Proxmox VM list
		vmList, err := pm.GetVMs(ctx, nodeData.Node)
		if err != nil {
			fmt.Printf("error getting proxmox vms: %s\n", err)
			os.Exit(1)
		}
		for _, vm := range vmList.Data {
			if vm.Status != "running" {
				continue
			}
			if excludedHosts.ContainsOne(vm.Name) {
				continue
			}
			hostname := fqdn(vm.Name)
			vmNames = append(vmNames, hostname)
			vmHosts[hostname] = hostInfo{node: nodeData.Node, vmid: vm.Vmid}
			tags := strings.Split(strings.Trim(vm.Tags, " "), ";")
			for _, tag := range tags {
				if tag == "" {
					continue
				}
				group := sanitizeGroupName(tag)
				if !slices.Contains(inv.All.Children, group) {
					inv.All.Children = append(inv.All.Children, group)
				}
				if _, exists := roles[group]; !exists {
					roles[group] = []string{}
				}
				roles[group] = append(roles[group], hostname)
			}
		}

		// Get Proxmox LXC list
		lxcs, err := pm.GetLxcs(ctx, nodeData.Node)
		if err != nil {
			fmt.Printf("error getting proxmox lxcs: %s\n", err)
			os.Exit(1)
		}
		for _, lxc := range lxcs.Data {
			if lxc.Status != "running" {
				continue
			}
			if excludedHosts.ContainsOne(lxc.Name) {
				continue
			}
			hostname := fqdn(lxc.Name)
			lxcNames = append(lxcNames, hostname)
			lxcHosts[hostname] = hostInfo{node: nodeData.Node, vmid: lxc.Vmid}
			tags := strings.Split(strings.Trim(lxc.Tags, " "), ";")
			for _, tag := range tags {
				if tag == "" {
					continue
				}
				group := sanitizeGroupName(tag)
				if !slices.Contains(inv.All.Children, group) {
					inv.All.Children = append(inv.All.Children, group)
				}
				if _, exists := roles[group]; !exists {
					roles[group] = []string{}
				}
				roles[group] = append(roles[group], hostname)
			}
		}
	}

	// Lookup IP addresses for ansible_host hostvars
	if Config.Proxmox.Lookup {
		for name, info := range lxcHosts {
			cfg, err := pm.GetLxcConfig(ctx, info.node, info.vmid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to get LXC config for %s: %v\n", name, err)
				continue
			}
			// Try Net0 through Net4
			for _, net := range []string{cfg.Data.Net0, cfg.Data.Net1, cfg.Data.Net2, cfg.Data.Net3, cfg.Data.Net4} {
				if ip := proxmox.ParseLxcIP(net); ip != "" {
					hostVarMap[name] = map[string]string{"ansible_host": ip}
					break
				}
			}
		}
		for name, info := range vmHosts {
			netResp, err := pm.GetQemuNetworkConfig(ctx, info.node, info.vmid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to get QEMU agent network info for %s: %v\n", name, err)
				continue
			}
			if ip := proxmox.FindQemuIPv4(netResp.Data.Result); ip != "" {
				hostVarMap[name] = map[string]string{"ansible_host": ip}
			}
		}
	}

	sort.Strings(inv.All.Children)
	sort.Strings(lxcNames)
	sort.Strings(vmNames)
	inv.Groups = make(ansible.InventoryGroupMap)
	inv.Groups["proxmox_lxcs"] = ansible.InventoryGroup{Hosts: lxcNames}
	inv.Groups["proxmox_vms"] = ansible.InventoryGroup{Hosts: vmNames}
	inv.Groups["ungrouped"] = ansible.InventoryGroup{Hosts: []string{}}

	keys := []string{}
	for k := range roles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sort.Strings(roles[k])
		if _, exists := inv.Groups[k]; !exists {
			inv.Groups[k] = ansible.InventoryGroup{Hosts: roles[k]}
		}
		if !slices.Contains(inv.All.Children, k) {
			inv.All.Children = append(inv.All.Children, k)
		}
	}

	// Handle --host: output hostvars for a single host
	if hostFlag != "" {
		vars := map[string]string{}
		if hv, ok := hostVarMap[hostFlag]; ok {
			vars = hv
		}
		str, err := json.MarshalIndent(vars, "", "   ")
		if err != nil {
			fmt.Printf("error marshalling json: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", str)
		os.Exit(0)
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

// fqdn returns the hostname with the configured domain appended, if set.
func fqdn(name string) string {
	if Config.Proxmox.Domain != "" {
		return name + "." + Config.Proxmox.Domain
	}
	return name
}

var groupNameRe = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// sanitizeGroupName converts a Proxmox tag to a valid Ansible group name.
// Ansible group names must match [a-zA-Z_][a-zA-Z0-9_]*.
func sanitizeGroupName(tag string) string {
	name := groupNameRe.ReplaceAllString(tag, "_")
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		name = "_" + name
	}
	return name
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
	viper.SetDefault("proxmox.domain", "")
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
