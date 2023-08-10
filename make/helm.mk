##@ ☸ Helm

KUSTOMIZE_ARGS ?= ""


.PHONY:helm-prepare
helm-prepare: manifests kustomize ## Prepare helm chart from kustomize resources
	$(KUSTOMIZE) build config/default -o helm/gko/templates/bundle.yaml 
	npx zx scripts/helm-transform.mjs
	$(MAKE) add-license

.PHONY: helm-template
helm-template: manifests kustomize helm-prepare ## Generates legacy bundle.yml file from helm chart
	helm template --include-crds  helm/gko -n gko-system > bundle.yml
