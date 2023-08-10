##@ 💻 Development

# Image URL to use to push to k3d local registry
K3D_IMG ?= k3d-graviteeio.docker.localhost:12345/gko
K3D_TAG ?= latest

K3DARGS ?= ""
.PHONY: k3d-init
k3d-init: ## Init and start the k3d cluster
	npx zx ./scripts/k3d.mjs $(K3DARGS)

.PHONY: k3d-start
k3d-start: ## Start the k3d cluster
	k3d cluster start graviteeio 

.PHONY: k3d-stop
k3d-stop: ## Stop the k3d cluster
	k3d cluster stop graviteeio

.PHONY: k3d-clean
k3d-clean: ## Delete the k3d cluster and docker registry
	k3d cluster delete graviteeio && k3d registry delete k3d-graviteeio.docker.localhost

.PHONY: k3d-gko-build
k3d-build: ## Build the controller image for k3d
	$(MAKE) docker-build IMG=$(K3D_IMG) TAG=$(K3D_TAG)

.PHONY: k3d-push
k3d-push: ## Push the controller image to the k3d registry
	$(MAKE) docker-push IMG=$(K3D_IMG) TAG=$(K3D_TAG)

.PHONY: k3d-deploy
k3d-deploy: ## Install operator helm chart to the k3d cluster
	$(MAKE) helm-prepare
	helm upgrade --install -n default --create-namespace gko helm/gko \
		--set manager.scope.cluster=false \
		--set manager.image.repository=$(K3D_IMG) \
		--set manager.image.tag=$(K3D_TAG)

.PHONY:
k3d-admin: ## Gain a kubernetes context with admin role on the k3d cluster
	kubectl config use-context k3d-graviteeio
	npx zx ./scripts/service-account.mjs

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: kustomize manifests helm-prepare ## Install CRDss into the current cluster
	kubectl apply -f helm/gko/crds

.PHONY: uninstall
uninstall: manifests kustomize  helm-prepare ## Uninstall CRDs from the current cluster
	kubectl delete -f helm/gko/crds

.PHONY: run
run: manifests generate ## Run a controller from your host
	go run ./main.go
