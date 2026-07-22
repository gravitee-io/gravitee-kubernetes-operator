##@ 🛠️ Tools

.PHONY: download
download: ## Download all project dependencies
	@echo "Downloading go.mod dependencies"
	@go mod download

.PHONY: install-tools
install-tools: download ## Installs non-Go tools (helm-unittest)
	@echo "Installing helm-unittest ..."
	@helm plugin list | grep -q unittest || helm plugin install https://github.com/quintush/helm-unittest > /dev/null 2>&1

.PHONY: setup-cursor
setup-cursor: ## Create .cursor/rules/ symlinks from .agent/rules/
	@mkdir -p .cursor/rules
	@for f in .agent/rules/*.md; do \
		ln -sf "../../$$f" ".cursor/rules/$$(basename $$f)"; \
	done
	@echo "Cursor rules symlinked from .agent/rules/"
