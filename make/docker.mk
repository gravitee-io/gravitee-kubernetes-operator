##@ 🐳 Docker

# Image URL to use all building/pushing image targets
IMG ?= controller:latest

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

