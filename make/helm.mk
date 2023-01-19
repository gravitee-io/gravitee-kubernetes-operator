##@ ☸ Helm

KUSTOMIZE_ARGS ?= ""


.PHONY:helm-prepare
helm-prepare: manifests kustomize ## generates a templated bundle.yaml file in helm chart
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default -o helm/gko/templates/bundle.yaml 
	npx zx scripts/helm-transform.mjs

.PHONY: helm-template
helm-template: manifests kustomize helm-prepare ## generates a templated bundle.yaml file in helm chart
	helm template --include-crds helm/gko -n gko-system > bundle.yaml
