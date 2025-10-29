##@ 🛠️ Tools

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

## Tool Binaries
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
GINKGO ?= $(LOCALBIN)/ginkgo
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
ADDLICENSE ?= $(LOCALBIN)/addlicense
CHAINSAW ?= $(LOCALBIN)/chainsaw
GOTESTSUM ?= $(LOCALBIN)/gotestsum
CLOUD_LB ?= $(LOCALBIN)/cloud-provider-kind

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
	@cd ./hack/tools && \
	for item in $$(find . -mindepth 1 -type d); do \
		pushd $${item} > /dev/null; \
		TOOL=$$(grep -e '^tool ' go.mod | sed -e s'/tool //'); \
		echo "Installing tool $${TOOL}"; \
		GOBIN=$(LOCALBIN) go install $${TOOL} & \
		popd > /dev/null; \
	done; \
	wait

.PHONY: install-tools
install-tools: install-go-tools ## Installs all required tools
	@echo "Installing helm-unittest ..."
	@helm plugin list | grep -q unittest || helm plugin install https://github.com/quintush/helm-unittest > /dev/null 2>&1

.PHONY: reinstall-tools
re-install-tools: clean-tools install-tools ## Clean and install tools again
