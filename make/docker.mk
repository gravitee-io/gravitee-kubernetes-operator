##@ 🐳 Docker

# Image URL to use all building/pushing image targets
IMG ?= graviteeio/kubernetes-operator
TAG ?= latest

# DEV

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG}:${TAG} .

.PHONY: docker-build-it
docker-build-it: ## Build the docker image with coverage info
	docker build -f Dockerfile.it -t ${IMG}:${TAG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}:${TAG}

# RELEASE

.PHONY: docker-build-release
docker-build-release: ## Build docker image with the manager.
	docker build -t ${IMG}:${TAG} -t ${IMG}:latest .

.PHONY: docker-push-release
docker-push-release: ## Push docker image with the manager.
	docker push  ${IMG}:${TAG}
	docker push ${IMG}:latest

