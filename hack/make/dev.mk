##@ 💻 Development

.PHONY: start-cluster
start-cluster: ## Init and start a local cluster
	@npx zx ./hack/scripts/run-kind.mjs


.PHONY: start-conformance-cluster
start-conformance-cluster: ## Init and start a local cluster for gateway-api conformance tests
	kind create cluster --config hack/kind/kind.conformance.yaml

.PHONY: delete-cluster
delete-cluster: ## Delete local cluster
	@kind delete cluster --name gravitee

.PHONY: use-cluster
use-cluster: ## Switch current kubectl context to local cluster
	@kubectl config use-context kind-gravitee

.PHONY: cluster-admin
cluster-admin: ## Gain a kubernetes context with admin role on the local cluster
	@kubectl config use-context kind-gravitee
	@npx zx ./hack/scripts/create-cluster-admin-sa.mjs
ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: ## Install CRDs into the current cluster
	@kubectl apply -f helm/gko/crds/gravitee.io

.PHONY: install-std-gateway-api
install-gateway-api: ## Install gateway-api CRDs into the current cluster
	@kubectl apply -f helm/gko/crds/kubernetes.io/gateway-api/standard

.PHONY: uninstall
uninstall: ## Uninstall CRDs from the current cluster
	@kubectl delete -f helm/gko/crds/gravitee.io

.PHONY: install-std-gateway-api
uninstall-gateway-api: ## Install gateway-api CRDs into the current cluster
	@kubectl delete -f helm/gko/crds/kubernetes.io/gateway-api/standard

.PHONY: run
run: ## Run a controller from your host
	@go run ./main.go
