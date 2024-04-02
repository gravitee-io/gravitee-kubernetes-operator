##@ 🧪 Test

.PHONY: helm-test
helm-test: helm-unittest
	@echo "Running helm unit tests ..."
	@helm unittest helm/gko

IT_ARGS ?= ""
COVER_PKG = "github.com/gravitee-io/gravitee-kubernetes-operator/..."
TIMEOUT = 380s 

.PHONY: test
it: use-cluster install ginkgo ## Run intgration tests
	$(GINKGO) $(IT_ARGS) --timeout $(TIMEOUT) --cover --coverprofile=cover.out --coverpkg=$(COVER_PKG)  test/integration/...

unit: ginkgo ## Run unit tests
	$(GINKGO) test/unit/...
