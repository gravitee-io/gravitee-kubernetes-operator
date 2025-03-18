##@ 🧪 Test

.PHONY: helm-test
helm-test: helm-unittest
	@echo "Running helm unit tests ..."
	@cd helm; helm unittest -f 'tests/**/*.yaml' gko


IT_ARGS ?= ""
TIMEOUT ?= 1200s 

.PHONY: it
it: use-cluster install ginkgo ## Run integration tests
	$(GINKGO) $(IT_ARGS) --timeout $(TIMEOUT)  test/integration/...

UT_ARGS ?= ""
unit: ginkgo ## Run unit tests
	$(GINKGO) $(UT_ARGS) test/unit/...
