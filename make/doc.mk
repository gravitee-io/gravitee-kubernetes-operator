##@ 📄 Documentation

.PHONY: reference
reference: crdoc ## Generate the reference documentation
	$(CRDOC) --resources config/crd/bases --output docs/api/reference.md


.PHONY: helm-reference ## generates helm chart documentation
helm-reference: helm-docs helm-prepare
	$(HELMDOCS) --dry-run
