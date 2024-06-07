##@ 🐳 Docker

# Image URL to use all building/pushing image targets
IMG ?= graviteeio/kubernetes-operator
TAG ?= latest

# DEV

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG}:${TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}:${TAG}


