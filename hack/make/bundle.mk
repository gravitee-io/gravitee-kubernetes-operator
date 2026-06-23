##@ 📦 OLM Bundle

VERSION ?= latest
BUNDLE_IMG ?= graviteeio/gko-bundle:$(VERSION)
CHANNELS ?= alpha
DEFAULT_CHANNEL ?= alpha

.PHONY: olm-bundle
olm-bundle: manifests ## Generate the OLM bundle from Helm chart and CSV base.
	npx zx hack/scripts/generate-bundle.mjs \
		--version $(VERSION) --img $(IMG) \
		--channels $(CHANNELS) --default-channel $(DEFAULT_CHANNEL)

.PHONY: olm-bundle-build
olm-bundle-build: ## Build the OLM bundle image.
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: olm-bundle-push
olm-bundle-push: ## Push the OLM bundle image.
	docker push $(BUNDLE_IMG)

OLM_TEST_VERSION ?= 0.0.0-test

.PHONY: olm-bundle-test
olm-bundle-test: ## Test OLM bundle install on a local Kind cluster.
	$(MAKE) olm-bundle olm-bundle-build VERSION=$(OLM_TEST_VERSION)
	$(MAKE) docker-build TAG=$(OLM_TEST_VERSION)
	npx zx hack/scripts/test-bundle.mjs \
		--version $(OLM_TEST_VERSION) --img $(IMG)
