##@ 🧪 Test

GOTESTARGS ?= ""
COVERPKG = "github.com/gravitee-io/gravitee-kubernetes-operator/..."
TIMEOUT = 380s 
.PHONY: test
test: manifests generate install ginkgo ## Run tests.
	kubectl config use-context k3d-graviteeio
	KUBEBUILDER_ASSETS=USE_EXISTING_CLUSTER=true $(GINKGO) $(GOTESTARGS) --timeout $(TIMEOUT) --cover --coverprofile=cover.out --coverpkg=$(COVERPKG)  ./...

.PHONY: helm-test
helm-test: helm-unittest
	@echo "Running helm unit tests ..."
	@helm unittest helm/gko
