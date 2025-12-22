# Proxmox Ansible Inventory

This program is used to provide a dynamic inventory to ansible for running LXC or Qemu virtual machines that are active in your proxmox cluster.
A cluster is not required, but if you have created a cluster with one or more Proxmox nodes defined, it will obtain a list of the running LXC containers
and Qemu virtual machines running across the cluster.

## Installation

1. Download the release from github for the Operating System and architecture of the environment you use to execute ansible. The list of available
operating systems and architectures are:

    * Darwin (macOS) - a universal binary is provided that will run on arm64 or x86_64 procesors
    * Linux - binaries are provided for arm64, i386, and x86_64 processors
    * Windows - binaries are provided for arm64, i386 and x86_64 processors

2. Optionally ceate an ansible hosts file to included static hosts in the inventory created by proxmox-ansible-inventory. This will produce a merged inventory of the
static hosts provided in the hosts file plus the LXC containers and Qemu virtual machines running in proxmox. For example:

    ##### hosts
    ```
    test1   ansible_host=192.168.1.20
    test2
    mac1    ansible_become=no ansible_user=bob 

3. Update your ansible.cfg to point to the location where you placed the proxmox-ansible-inventory executable. The example below shows that I've
placed the executable in the same directory as the ansible.cfg file.

    ##### ansible.cfg
    ```
    [defaults]
    collections_scan_sys_path = false
    cow_selection = tux
    deprecation_warnings = false
    gathering = smart
    gather_subset = !hardware
    gather_timeout = 10
    inventory = ./proxmox-ansible-inventory
    interpreter_python = auto_silent
    timeout=60
    verbosity=1
    [ssh_connection]
    scp_if_ssh=True
    transfer_method=scp

4. Run the "_ansible-inventory --list_" command, which should produce output similar to the following:

    ```
    {
        "_meta": {
            "hostvars": {
                "test1": {
                    "ansible_host": "192.168.1.20"
                },
                "mac1": {
                    "ansible_become": "no",
                    "ansible_user": "bob"
                }
            }
        },
        "all": {
            "children": [
                "ungrouped",
                "proxmox_lxcs",
                "proxmox_vms",
                "static"
            ]
        },
        "proxmox_lxcs": {
            "hosts": [
                "docker1",
                "debian1",
                "network-test",
                "netbox1",
                "postgres1",
                "redis1",
                "docker2"
            ]
        },
        "proxmox_vms": {
            "hosts": [
                "pbs1",
                "debian2",
                "pbs2"
            ]
        },
        "static": {
            "hosts": [
                "test1",
                "test2",
                "mac1"
            ]
        }
    }
