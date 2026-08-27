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

CONFORMANCE_TIMEOUT ?= 30m
CONFORMANCE_JUNIT_FILE ?= /tmp/junit/reports/conformance.xml
CONFORMANCE_RERUN ?= 2
CONFORMANCE_RERUN_REPORT ?= $(dir $(CONFORMANCE_JUNIT_FILE))conformance-reruns.txt

# Reruns cost nothing on a green run, and roughly the failed test alone on a red one, since
# CleanupBaseResources is false and a rerun reuses the base gateways rather than rebuilding
# them. They mitigate flakiness, they do not fix it, so two guards keep them honest:
# max-failures declines to rerun anything when the first pass failed broadly, which is
# breakage rather than flakiness, and the report records what was retried so the flake rate
# stays measurable instead of vanishing into a green tick. Publish that report from CI.
#
# Set CONFORMANCE_RERUN=0 when a run is meant to be evidence, in particular the consecutive
# runs required by .agent/rules/gateway-standards.md: a retried pass is not a pass there.
ifneq ($(CONFORMANCE_RERUN),0)
CONFORMANCE_RERUN_FLAGS := --rerun-fails=$(CONFORMANCE_RERUN) --rerun-fails-max-failures=5 --rerun-fails-report=$(CONFORMANCE_RERUN_REPORT)
endif

.PHONY: conformance
conformance: ## Run conformance tests
	@mkdir -p $(dir $(CONFORMANCE_JUNIT_FILE))
	GATEWAY_API_MATCH_ACROSS_ROUTES=true go tool gotestsum --format=testname --junitfile $(CONFORMANCE_JUNIT_FILE) $(CONFORMANCE_RERUN_FLAGS) --packages="./test/conformance/kubernetes.io/gateway-api/standard/..." -- -timeout $(CONFORMANCE_TIMEOUT) -args --gateway-class=gravitee-gateway
