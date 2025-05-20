##@ 🧪 Test

.PHONY: helm-test
helm-test:
	@echo "Running helm unit tests ..."
	@cd helm; helm unittest -f 'tests/**/*.yaml' gko


IT_ARGS ?= ""
TIMEOUT ?= 1200s 

.PHONY: it
it: use-cluster install ## Run integration tests
	$(GINKGO) $(IT_ARGS) --timeout $(TIMEOUT)  test/integration/...

UT_ARGS ?= ""
unit:  ## Run unit tests
	$(GINKGO) $(UT_ARGS) test/unit/...
