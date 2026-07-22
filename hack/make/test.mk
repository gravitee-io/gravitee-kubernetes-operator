##@ 🧪 Test

.PHONY: helm-test
helm-test:
	@echo "Running helm unit tests ..."
	@cd helm; helm unittest -f 'tests/**/*.yaml' gko


IT_ARGS ?= ""
TIMEOUT ?= 1200s 

.PHONY: it
it: use-cluster install ## Run integration tests
	go tool ginkgo $(IT_ARGS) --timeout $(TIMEOUT)  test/integration/...

UT_ARGS ?= ""
.PHONY: unit
unit:  ## Run unit tests
	go tool ginkgo $(UT_ARGS) test/unit/...

.PHONY: e2e
e2e:  ## Run all end to end tests
	$(CHAINSAW) test --config test/e2e/chainsaw/config.yaml

.PHONY: e2e-focus
e2e-focus:  ## Run end to end tests and focus on test having the focus label set to true
	$(CHAINSAW) test --config test/e2e/chainsaw/config.yaml --selector focus=true

.PHONY: conformance
conformance: ## Run conformance tests
	go tool gotestsum --format=testname --packages="./test/conformance/kubernetes.io/gateway-api/standard/..." -- -args --gateway-class=gravitee-gateway
