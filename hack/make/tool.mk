##@ 🛠️ Tools

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
GINKGO ?= $(LOCALBIN)/ginkgo
CRDOC ?= $(LOCALBIN)/crdoc
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
ADDLICENSE ?= $(LOCALBIN)/addlicense

.PHONY: clean-tools
clean-tools: $(LOCALBIN) ## Cleans (delete) all binary tools
	@echo "Cleaning tools"
	@find $(LOCALBIN) -type f -delete

.PHONY: download
download: ## Download all project dependencies
	@echo "Downloading go.mod dependencies"
	@go mod download

.PHONY: install-go-tools
install-go-tools: download ## Installs all required GO tools
	@echo "Installing GO tools"
	@cat hack/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -I % sh -c 'GOBIN=$(LOCALBIN) go install %'

.PHONY: install-tools
install-tools: install-go-tools ## Installs all required tools
	@echo "Installing helm-unittest ..."
	@helm plugin list | grep -q unittest || helm plugin install https://github.com/quintush/helm-unittest > /dev/null 2>&1

.PHONY: reinstall-tools
re-install-tools: clean-tools install-tools ## Clean and install tools again
