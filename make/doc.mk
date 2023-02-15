##@ 📄 Documentation

.PHONY: reference
reference: crdoc ## Generate the CRDs reference documentation
	$(CRDOC) --resources config/crd/bases --output docs/api/reference.md


.PHONY: helm-reference
helm-reference: helm-docs helm-prepare ## Generates helm chart documentation
	@$(HELMDOCS) --chart-search-root=helm --dry-run > helm/gko/README.md
