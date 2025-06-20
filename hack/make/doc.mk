##@ 📄 Documentation

.PHONY: reference
reference: manifests-for-docs ## Generate the CRDs reference documentation
	$(CRDOC) --resources docs/api/crd --output docs/api/reference.md --template hack/crdoc/markdown.tmpl -c hack/crdoc/toc.yaml
	@rm -rf docs/api/crd

.PHONY:
helm-reference: ## Generates helm chart documentation
	npx @bitnami/readme-generator-for-helm -v helm/gko/values.yaml -r helm/gko/README.md
