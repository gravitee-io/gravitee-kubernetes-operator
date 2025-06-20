##@ 📄 Documentation

.PHONY: reference
reference:
	@bin/crd-ref-docs \
		--max-depth=100 \
		--source-path=${PWD}/api \
		--config=.crd-ref-docs.yaml \
		--renderer=markdown \
		--output-path=${PWD}/docs/api/reference.md
	@npx zx hack/scripts/clean-reference.mjs

.PHONY:
helm-reference: ## Generates helm chart documentation
	npx @bitnami/readme-generator-for-helm -v helm/gko/values.yaml -r helm/gko/README.md
