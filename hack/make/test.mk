##@ 🧪 Test

.PHONY: helm-test
helm-test:
	@echo "Running helm unit tests ..."
	@cd helm; helm unittest -f 'tests/**/*.yaml' gko


IT_ARGS ?= ""
TIMEOUT ?= 1200s 

.PHONY: it
it: use-cluster install install-go-tools ## Run integration tests
	$(GINKGO) $(IT_ARGS) --timeout $(TIMEOUT)  test/integration/...

UT_ARGS ?= ""
.PHONY: unit
unit:  ## Run unit tests
	$(GINKGO) $(UT_ARGS) test/unit/...

.PHONY: e2e
e2e:  ## Run all end to end tests (Playwright)
	npm --prefix test/platform-test run e2e

.PHONY: conformance
conformance: install-go-tools ## Run end to end tests and focus on test having the focus label set to true
	$(GOTESTSUM) --format=testname --packages="./test/conformance/kubernetes.io/gateway-api/standard/..." -- -args --gateway-class=gravitee-gateway
