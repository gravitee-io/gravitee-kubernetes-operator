##@ 📄 Documentation

.PHONY: reference
reference: ## Generate the CRDs reference documentation
	$(CRDOC) --resources helm/gko/crds --output docs/api/reference.md --template hack/crdoc/markdown.tmpl -c hack/crdoc/toc.yaml

.PHONY:
helm-reference: ## Generates helm chart documentation
	npx @bitnami/readme-generator-for-helm -v helm/gko/values.yaml -r helm/gko/README.md
