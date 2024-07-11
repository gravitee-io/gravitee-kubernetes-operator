##@ 🧪 Test

.PHONY: helm-test
helm-test: helm-unittest
	@echo "Running helm unit tests ..."
	@helm unittest helm/gko

IT_ARGS ?= ""
TIMEOUT ?= 1200s 

.PHONY: test
it: use-cluster install ginkgo ## Run intgration tests
	$(GINKGO) $(IT_ARGS) --timeout $(TIMEOUT)  test/integration/...

UT_ARGS ?= ""
unit: ginkgo ## Run unit tests
	$(GINKGO) $(UT_ARGS) test/unit/...
