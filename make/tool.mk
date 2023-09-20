##@ 🛠️ Tools

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
GINKGO ?= $(LOCALBIN)/ginkgo
CRDOC ?= $(LOCALBIN)/crdoc
GOLANGCILINT ?= $(LOCALBIN)/golangci-lint
ADDLICENSE ?= $(LOCALBIN)/addlicense
HELMDOCS ?= $(LOCALBIN)/helm-docs

ALL_TOOLS = kustomize controller-gen envtest ginkgo crdoc golangci-lint addlicense helm-docs helm-unittest

## Tool Versions
KUSTOMIZE_VERSION ?= v4.5.7
CONTROLLER_TOOLS_VERSION ?= v0.11.1
ENVTEST_K8S_VERSION = 1.24

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	@echo "Installing kustomize ..."
	@test -f $(LOCALBIN)/kustomize || curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN)

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	@echo "Installing controller-gen .."
	@GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	@echo "Installing envtest ..."
	@GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo cli locally if necessary.
$(GINKGO): $(LOCALBIN)
	@echo "Installing ginkgo ..."
	@GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@latest
 
.PHONY: crdoc ## Download crdoc cli locally if necessary.
crdoc: $(CRDOC)
$(CRDOC): $(LOCALBIN)
	@echo "Installing crdoc ..."
	@GOBIN=$(LOCALBIN) go install fybrik.io/crdoc@latest

.PHONY: golangci-lint ## Download golangci-lint cli locally if necessary
golangci-lint: $(GOLANGCILINT)
$(GOLANGCILINT): $(LOCALBIN)
	@echo "Installing golangci-lint ..."
	@GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: addlicense
addlicense: $(ADDLICENSE) ## Download addlicense cli locally if necessary.
$(ADDLICENSE): $(LOCALBIN)
	@echo "Installing addlicense ..."
	@GOBIN=$(LOCALBIN) go install github.com/google/addlicense@latest

.PHONY: helm-docs
helm-docs: $(HELMDOCS) ## Download helmdocs cli locally if necessary.
$(HELMDOCS): $(LOCALBIN)
	@echo "Installing helm-docs ..."
	@GOBIN=$(LOCALBIN) go install github.com/norwoodj/helm-docs/cmd/helm-docs@latest

.PHONY: helm-unittest
helm-unittest: ## Install helm-unittest plugin if necessary.
	@echo "Installing helm-unittest ..."
	@helm plugin list | grep -q unittest || helm plugin install https://github.com/quintush/helm-unittest > /dev/null 2>&1

.PHONY: all-tools
install-tools: clean-tools $(ALL_TOOLS) ## Install all binary tools (use -j to run in parallel)

.PHONY: clean-tools 
clean-tools: ## Clean (delete) all binary tools
	@echo "Cleaning tools"
	@find $(LOCALBIN) -type f -delete 
