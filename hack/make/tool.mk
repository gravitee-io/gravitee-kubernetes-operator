##@ 🛠️ Tools

.PHONY: download
download: ## Download all project dependencies
	@echo "Downloading go.mod dependencies"
	@go mod download

.PHONY: install-tools
install-tools: download ## Installs non-Go tools (helm-unittest)
	@echo "Installing helm-unittest ..."
	@helm plugin list | grep -q unittest || helm plugin install https://github.com/quintush/helm-unittest > /dev/null 2>&1
