##@ 📄 Documentation

.PHONY: reference
<<<<<<< HEAD
reference: ## Generate the CRDs reference documentation
	$(CRDOC) --resources helm/gko/crds --output docs/api/reference.md --template hack/crdoc/markdown.tmpl -c hack/crdoc/toc.yaml
=======
reference:
	@bin/crd-ref-docs \
		--max-depth=100 \
		--source-path=${PWD}/api \
		--config=.crd-ref-docs.yaml \
		--renderer=markdown \
		--output-path=${PWD}/docs/api/reference.md
	@npx zx hack/scripts/clean-reference.mjs
>>>>>>> 56d297b (docs: generate docs using crd-ref-docs)

.PHONY:
helm-reference: ## Generates helm chart documentation
	npx @bitnami/readme-generator-for-helm -v helm/gko/values.yaml -r helm/gko/README.md
