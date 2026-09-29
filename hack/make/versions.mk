##@ 📌 Versions

# Kubernetes version of the APIM kind cluster (make start-cluster) and of the envtest control plane.
# KIND_VERSION must ship a kindest/node image for it: see the kind release notes.
# The Gateway API conformance cluster keeps its own versions (.circleci/config.yml).
K8S_VERSION ?= 1.37.0
KIND_VERSION ?= v0.33.0

.PHONY: print-k8s-version
print-k8s-version: ## Print the Kubernetes version of the kind cluster and envtest
	@echo $(K8S_VERSION)

.PHONY: print-kind-version
print-kind-version: ## Print the kind version that creates the kind cluster
	@echo $(KIND_VERSION)
