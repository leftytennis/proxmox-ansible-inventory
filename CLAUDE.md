# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go CLI tool that queries the Proxmox VE API to generate a dynamic Ansible inventory. It discovers all running LXC containers and Qemu VMs across Proxmox cluster nodes and outputs JSON in Ansible's dynamic inventory format. Proxmox tags on VMs/LXCs become Ansible groups.

## Build & Run

```bash
# Build
go build -o proxmox-ansible-inventory .

# Build with version info (matches goreleaser ldflags)
go build -ldflags "-X main.GitVersion=dev -X main.GitSha=$(git rev-parse HEAD) -X main.GitDate=$(date -u +%Y-%m-%d)" -o proxmox-ansible-inventory .

# Run (requires config file, see below)
./proxmox-ansible-inventory --list

# Release build (uses goreleaser)
goreleaser release --snapshot --clean

# Vendor dependencies
go mod tidy && go mod vendor
```

There are no tests yet.

## Configuration

Config file: `.proxmox-ansible-inventory.yml` — searched in `.` and `$HOME/.config/proxmox-ansible-inventory/`. Environment variables use prefix `PAI_`. See `.proxmox-ansible-inventory.yml.example` for format.

Required config values: `proxmox.api.user`, `proxmox.api.token`, `proxmox.api.secret`, `proxmox.api.url`.

## Architecture

Three packages with a simple data flow: **config** → **proxmox** → **ansible** → JSON output.

- **`main.go`** — Entry point. Sets up viper config, iterates Proxmox nodes, collects VMs/LXCs, maps Proxmox tags to Ansible groups, and outputs the inventory as JSON. Version info (`GitVersion`, `GitSha`, `GitDate`) is injected via ldflags at build time.
- **`config/`** — Config struct (`Params`) with viper `mapstructure` tags. `CheckRequiredValues()` validates API credentials are present.
- **`proxmox/`** — HTTP client for the Proxmox VE REST API (`/api2/json`). Authenticates via PVE API token header. Key methods: `GetNodes`, `GetVMs`, `GetLxcs`. Also has methods for fetching individual configs and QEMU agent network info (`GetLxcConfig`, `GetVMConfig`, `GetQemuNetworkConfig`).
- **`ansible/`** — Inventory types matching Ansible's JSON inventory format. Custom `MarshalJSON` on `Inventory` merges the fixed fields (`_meta`, `all`) with the dynamic `Groups` map into a single flat JSON object.

## Key Details

- Dependencies are vendored (`vendor/`). Use `go mod vendor` after changing deps.
- Proxmox tags (semicolon-delimited in the API) become Ansible group names. Each tag creates a group with its tagged hosts.
- Only **running** VMs/LXCs are included. Hosts in `proxmox.exclude` config list are filtered out using `mapset`.
- The `Inventory.Groups` field uses `json:"-"` and is merged into the top-level JSON via the custom `MarshalJSON` method — this is important to understand when modifying the output format.
