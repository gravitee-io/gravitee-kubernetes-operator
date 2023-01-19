##@ Test

GOTESTARGS ?= ""
.PHONY: test
test: manifests generate install ginkgo ## Run tests.
	kubectl config use-context k3d-graviteeio
	KUBEBUILDER_ASSETS=USE_EXISTING_CLUSTER=true $(GINKGO) $(GOTESTARGS) --timeout 380s --cover --coverprofile=cover.out ./...
