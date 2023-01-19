##@ 📄 Documentation

.PHONY: reference
reference: crdoc ## Generate the reference documentation
	$(CRDOC) --resources config/crd/bases --output docs/api/reference.md
