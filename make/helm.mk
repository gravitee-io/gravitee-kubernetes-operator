##@ ☸ Helm

.PHONY: helm-template
helm-template: manifests ## Generates legacy bundle.yml file from helm chart
	helm template --include-crds  helm/gko -n gko-system > bundle.yml
