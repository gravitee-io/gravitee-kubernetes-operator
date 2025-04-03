##@ 🧹 Lint

ALL_LINT=lint-commits lint-sources lint-licenses

.PHONY: lint
lint: lint-commits lint-sources lint-licenses ## Run all linters and fail on error

.PHONY: lint-commits
lint-commits:  ## Run commitlint and fail on error
	@echo "Linting commits ..."
	@npm i -g @commitlint/config-conventional @commitlint/cli
	@commitlint -x @commitlint/config-conventional --edit

.PHONY: lint-sources
lint-sources: ## Run golangci-lint and fail on error
	@echo "Linting go sources ..."
	@$(GOLANGCILINT) --concurrency 2 run ./...
	@npx --yes prettier --check hack/scripts

.PHONY: lint-licenses
lint-licenses: ## Run addlicense linter and fail on error
	@echo "Checking license headers ..."
	@$(ADDLICENSE) -check -f LICENSE_TEMPLATE.txt \
		-ignore ".circleci/**" \
		-ignore ".mergify.yml" \
		-ignore "config/**" \
		-ignore ".idea/**" .

.PHONY: add-license
add-license: ## Add license headers to files
	@echo "Adding license headers ..."
	@$(ADDLICENSE) -f LICENSE_TEMPLATE.txt \
		-ignore ".circleci/**" \
		-ignore ".mergify.yml" \
		-ignore "config/**" \
		-ignore ".idea/**" .

.PHONY: clean-tools ## Run all linters
lint: $(ALL_LINT)

.PHONY: lint-fix
lint-fix: ## Fix whatever golangci-lint can fix and add licenses headers
	@$(GOLANGCILINT) run ./... --fix
	@$(MAKE) add-license
	@npx --yes prettier --write hack/scripts
