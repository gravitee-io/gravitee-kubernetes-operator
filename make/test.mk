##@ 🧪 Test

COVERPKG = "github.com/gravitee-io/gravitee-kubernetes-operator/..."
TIMEOUT = 600s 
.PHONY: test
it: ## Run integragtion tests.
	KUBEBUILDER_ASSETS=USE_EXISTING_CLUSTER=true $(GINKGO) -procs=4 --timeout $(TIMEOUT) --cover --coverprofile=cover.out --coverpkg=$(COVERPKG)  ./test/integration/...

.PHONY: helm-test
helm-test: helm-unittest
	@echo "Running helm unit tests ..."
	@helm unittest helm/gko
