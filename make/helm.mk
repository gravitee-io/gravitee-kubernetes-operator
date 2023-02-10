##@ ☸ Helm

KUSTOMIZE_ARGS ?= ""


.PHONY:helm-prepare
helm-prepare: manifests kustomize ## prepare helm chart from kustomize resources
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default -o helm/gko/templates/bundle.yaml 
	npx zx scripts/helm-transform.mjs

.PHONY: helm-template
helm-template: manifests kustomize helm-prepare ## generates a templated (legacy) bundle.yml
	helm template --include-crds  helm/gko -n gko-system > bundle.yml

.PHONY: helm-document ## generates helm chart documentation
helm-document: helm-docs
	$(HELMDOCS) --dry-run
