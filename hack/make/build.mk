##@ 🔨Build

.PHONY: build
build: generate ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: manifests
manifests: ## Generate CustomResourceDefinition objects.
	go tool controller-gen crd:maxDescLen=100 paths="./api/..." output:crd:artifacts:config=crds/gravitee.io
	@npx zx hack/scripts/annotate-crds.mjs
	$(MAKE) add-license

.PHONY: manifests-for-docs
manifests-for-docs: ## Generate CustomResourceDefinition objects.
	@mkdir -p docs/api/crd
	go tool controller-gen crd paths="./api/..." output:crd:artifacts:config=docs/api/crd

.PHONY: generate
generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	go run ./hack/crdgen -module github.com/gravitee-io/gravitee-automation-sdk/am-sdk/v2 -spec openapi/openapi.yaml -overlay hack/crdgen/am/overlay.yaml -out hack/crdgen/am/openapi.gen.yaml
	go generate ./api/model/am/...
	go tool controller-gen object:headerFile="hack/license.go.txt" paths="./..."

