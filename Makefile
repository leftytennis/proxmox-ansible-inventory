BINARY_NAME := proxmox-ansible-inventory
INSTALL_DIR := $(HOME)/develop/ansible
LDFLAGS := -X main.GitVersion=dev -X main.GitSha=$$(git rev-parse HEAD) -X main.GitDate=$$(date -u +%Y-%m-%d)

.PHONY: build install-local release-local clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

install-local:
	go build -ldflags "$(LDFLAGS)" -o $(INSTALL_DIR)/$(BINARY_NAME) .

release-local:
	goreleaser release --snapshot --clean
	cp dist/$(BINARY_NAME)_darwin_all/$(BINARY_NAME) $(INSTALL_DIR)/

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
