<<<<<<< HEAD
=======
##@ 🔨Build

.PHONY: build
build: generate ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: manifests
manifests: controller-gen ## Generate CustomResourceDefinition objects.
	$(CONTROLLER_GEN) crd:maxDescLen=100 paths="./api/..." output:crd:artifacts:config=helm/gko/crds/gravitee.io
	@npx zx scripts/annotate-crds.mjs
	$(MAKE) add-license

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

>>>>>>> a1e5e48b (feat: add controller for gateway class)
