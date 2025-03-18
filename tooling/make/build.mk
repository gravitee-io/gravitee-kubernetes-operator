##@ 🔨Build

.PHONY: build
build: generate ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: manifests
manifests: controller-gen ## Generate CustomResourceDefinition objects.
	$(CONTROLLER_GEN) crd paths="./..." output:crd:artifacts:config=helm/gko/crds
	@npx zx tooling/scripts/annotate-crds.mjs
	$(MAKE) add-license

.PHONY: generate
generate: controller-gen ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="tooling/license.go.txt" paths="./..."

