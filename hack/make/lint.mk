##@ 🧹 Lint

.PHONY: lint
lint: lint-sources lint-licenses ## Run all linters and fail on error

.PHONY: lint-sources
lint-sources: lint-vet lint-staticcheck lint-revive lint-prettier ## Run all source linters in parallel (use make -j)

.PHONY: lint-vet
lint-vet: ## Run go vet
	@echo "Running go vet ..."
	@go vet ./...

.PHONY: lint-staticcheck
lint-staticcheck: ## Run staticcheck
	@echo "Running staticcheck ..."
	@go tool staticcheck ./...

.PHONY: lint-revive
lint-revive: ## Run revive
	@echo "Running revive ..."
	@go tool revive -config .revive.toml -formatter friendly ./...

.PHONY: lint-prettier
lint-prettier: ## Run prettier on scripts
	@npx --yes prettier --check hack/scripts

.PHONY: lint-commits
lint-commits: ## Run commitlint and fail on error
	@echo "Linting commits ..."
	@npm i -g @commitlint/config-conventional @commitlint/cli
	@commitlint -x @commitlint/config-conventional --edit

.PHONY: lint-licenses
lint-licenses: ## Run addlicense linter and fail on error
	@echo "Checking license headers ..."
	@go tool addlicense -check -f LICENSE_TEMPLATE.txt \
		-ignore ".circleci/**" \
		-ignore ".mergify.yml" \
		-ignore ".crd-ref-docs.yaml" \
		-ignore ".idea/**" \
		-ignore "crds/kubernetes.io/**" \
		-ignore "examples/gateway-api/**" \
		-ignore "**/kubernetes.io/gateway-api/report/**" \
		-ignore "examples/gateway-api/**" \
		. 

.PHONY: add-license
add-license: ## Add license headers to files
	@echo "Adding license headers ..."
	@go tool addlicense -f LICENSE_TEMPLATE.txt \
		-ignore ".circleci/**" \
		-ignore ".mergify.yml" \
		-ignore ".crd-ref-docs.yaml" \
		-ignore ".idea/**" \
		-ignore "crds/kubernetes.io/**" \
		-ignore "**/kubernetes.io/gateway-api/report/**" \
		-ignore "examples/gateway-api/**" \
		. 

.PHONY: lint-fix
lint-fix: ## Auto-fix what can be fixed and add license headers
	@find . -name '*.go' -not -name 'zz_generated*' -not -path './vendor/*' | xargs go tool goimports -w
	@$(MAKE) add-license
	@npx --yes prettier --write hack/scripts
