# Ensure Make is run with bash shell as some syntax below is bash-specific
SHELL:=/usr/bin/env bash

# Define registries
STAGING_REGISTRY ?= ghcr.io/cohesity

IMAGE_NAME ?= cluster-api-byoh-controller
TAG ?= dev
RELEASE_DIR := _dist

# Image URL to use all building/pushing image targets
IMG ?= ${STAGING_REGISTRY}/${IMAGE_NAME}:${TAG}
BYOH_BASE_IMG = byoh/node:e2e
BYOH_BASE_IMG_DEV = byoh/node:dev
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)

REPO_ROOT := $(shell pwd)
GINKGO_FOCUS  ?=
GINKGO_SKIP ?=
GINKGO_NODES  ?= 1
E2E_CONF_FILE  ?= ${REPO_ROOT}/test/e2e/config/provider.yaml
ARTIFACTS ?= ${REPO_ROOT}/_artifacts
SKIP_RESOURCE_CLEANUP ?= false
USE_EXISTING_CLUSTER ?= false
EXISTING_CLUSTER_KUBECONFIG_PATH ?=
GINKGO_NOCOLOR ?= false

TOOLS_DIR := $(REPO_ROOT)/hack/tools
BIN_DIR := bin
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)
GINKGO_PKG := github.com/onsi/ginkgo/v2/ginkgo

BYOH_TEMPLATES := $(REPO_ROOT)/test/e2e/data/infrastructure-provider-byoh

LDFLAGS := -w -s $(shell hack/version.sh)
STATIC=-extldflags '-static'

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.DEFAULT_GOAL := help

.PHONY: all
all: build

HOST_AGENT_DIR ?= agent

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests
manifests: ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_bootstrapkubeconfigs.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_byoclusters.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_byoclustertemplates.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_byohosts.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_byomachines.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_byomachinetemplates.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_k8sinstallerconfigs.yaml
	$(YQ) -i eval 'del(.metadata.creationTimestamp)' config/crd/bases/infrastructure.cluster.x-k8s.io_k8sinstallerconfigtemplates.yaml

.PHONY: generate
generate: ## Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	GOOS=linux go vet ./...

.PHONY: test
test: manifests generate fmt vet setup-envtest test-coverage ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

test-coverage: prepare-byoh-docker-host-image ## Run test-coverage
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; $(GINKGO) --randomize-all -r --cover --coverprofile=cover.out --output-dir=. --skip-package=test .

agent-test: prepare-byoh-docker-host-image ## Run agent tests
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; $(GINKGO) --randomize-all -r $(HOST_AGENT_DIR) --coverprofile cover.out

controller-test: ## Run controller tests
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; $(GINKGO) --randomize-all controllers/infrastructure --coverprofile cover.out --vv

webhook-test: ## Run webhook tests
	source ./scripts/fetch_ext_bins.sh; fetch_tools; setup_envs; $(GINKGO) internal/webhook/infrastructure/v1beta1 --coverprofile cover.out

# test-e2e: take-user-input docker-build prepare-byoh-docker-host-image $(GINKGO) cluster-templates-e2e ## Run the end-to-end tests
# 	$(GINKGO) -v -trace -tags=e2e -focus="$(GINKGO_FOCUS)" $(_SKIP_ARGS) -nodes=$(GINKGO_NODES) --noColor=$(GINKGO_NOCOLOR) $(GINKGO_ARGS) test/e2e -- \
# 	    -e2e.artifacts-folder="$(ARTIFACTS)" \
# 	    -e2e.config="$(E2E_CONF_FILE)" \
# 	    -e2e.skip-resource-cleanup=$(SKIP_RESOURCE_CLEANUP) -e2e.use-existing-cluster=$(USE_EXISTING_CLUSTER) \
# 		-e2e.existing-cluster-kubeconfig-path=$(EXISTING_CLUSTER_KUBECONFIG_PATH)

# TODO(user): To use a different vendor for e2e tests, modify the setup under 'tests/e2e'.
# The default setup assumes Kind is pre-installed and builds/loads the Manager Docker image locally.
# CertManager is installed by default; skip with:
# - CERT_MANAGER_INSTALL_SKIP=true
KIND_CLUSTER ?= byoh-test-e2e

.PHONY: setup-test-e2e
setup-test-e2e: ## Set up a Kind cluster for e2e tests if it does not exist
	@command -v $(KIND) >/dev/null 2>&1 || { \
		echo "Kind is not installed. Please install Kind manually."; \
		exit 1; \
	}
	@case "$$($(KIND) get clusters)" in \
		*"$(KIND_CLUSTER)"*) \
			echo "Kind cluster '$(KIND_CLUSTER)' already exists. Skipping creation." ;; \
		*) \
			echo "Creating Kind cluster '$(KIND_CLUSTER)'..."; \
			$(KIND) create cluster --name $(KIND_CLUSTER) ;; \
	esac

.PHONY: test-e2e
test-e2e: setup-test-e2e manifests generate fmt vet ## Run the e2e tests. Expected an isolated environment using Kind.
	KIND=$(KIND) KIND_CLUSTER=$(KIND_CLUSTER) go test -tags=e2e ./test/e2e/ -v -ginkgo.v
	$(MAKE) cleanup-test-e2e

.PHONY: cleanup-test-e2e
cleanup-test-e2e: ## Tear down the Kind cluster used for e2e tests
	@$(KIND) delete cluster --name $(KIND_CLUSTER)

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	GOOS=linux $(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	GOOS=linux $(GOLANGCI_LINT) run --fix

.PHONY: lint-config
lint-config: golangci-lint ## Verify golangci-lint linter configuration
	$(GOLANGCI_LINT) config verify

##@ Build

.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/main.go

# If you wish to build the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	$(CONTAINER_TOOL) build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	$(CONTAINER_TOOL) push ${IMG}

prepare-byoh-docker-host-image:
	docker build test/e2e -f test/e2e/BYOHDockerFile -t ${BYOH_BASE_IMG}

prepare-byoh-docker-host-image-dev:
	docker build test/e2e -f docs/BYOHDockerFileDev -t ${BYOH_BASE_IMG_DEV}

cluster-templates-v1beta1: ## Generate cluster templates for v1beta1
	$(KUSTOMIZE) build $(BYOH_TEMPLATES)/v1beta1/templates/vm --load-restrictor LoadRestrictionsNone > $(BYOH_TEMPLATES)/v1beta1/templates/vm/cluster-template.yaml
	$(KUSTOMIZE) build $(BYOH_TEMPLATES)/v1beta1/templates/docker --load-restrictor LoadRestrictionsNone > $(BYOH_TEMPLATES)/v1beta1/templates/docker/cluster-template.yaml

##@ Test

cluster-templates: cluster-templates-v1beta1

cluster-templates-e2e:
	$(KUSTOMIZE) build $(BYOH_TEMPLATES)/v1beta1/templates/e2e --load-restrictor LoadRestrictionsNone > $(BYOH_TEMPLATES)/v1beta1/templates/e2e/cluster-template.yaml

define WARNING
#####################################################################################################

** WARNING **
These tests modify system settings - and do **NOT** revert them at the end of the test.
A list of changes can be found below. We **highly** recommend running these tests in a VM.

Running e2e tests locally will change the following host config
- enable the kernel modules: overlay & bridge network filter
- create a systemwide config that will enable those modules at boot time
- enable ipv4 & ipv6 forwarding
- create a systemwide config that will enable the forwarding at boot time
- reload the sysctl with the applied config changes so the changes can take effect without restarting
- disable unattended OS updates

#####################################################################################################
endef
export WARNING

.PHONY: take-user-input
take-user-input:
	@echo "$$WARNING"
	@read -p "Do you want to proceed [Y/n]?" REPLY; \
	if [[ $$REPLY = "Y" || $$REPLY = "y" ]]; then echo starting e2e test; exit 0 ; else echo aborting; exit 1; fi

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name byoh-builder
	$(CONTAINER_TOOL) buildx use byoh-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm byoh-builder
	rm Dockerfile.cross

.PHONY: build-installer
build-installer: manifests generate ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > dist/install.yaml

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: manifests ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) apply -f -

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/crd | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: manifests ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image ghcr.io/cohesity/cluster-api-provider-bringyourownhost-controller=${IMG}
	$(KUSTOMIZE) build config/default | $(KUBECTL) apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	$(KUSTOMIZE) build config/default | $(KUBECTL) delete --ignore-not-found=$(ignore-not-found) -f -

publish-infra-yaml: # Generate infrastructure-components.yaml for the provider
	cd config/manager && $(KUSTOMIZE) edit set image ghcr.io/cohesity/cluster-api-provider-bringyourownhost-controller=${IMG}
	$(KUSTOMIZE) build config/default > infrastructure-components.yaml

host-agent-binaries: ## Builds the binaries for the host-agent
	RELEASE_BINARY=./byoh-hostagent GOOS=linux GOARCH=amd64 GOLDFLAGS="$(LDFLAGS) $(STATIC)" \
	HOST_AGENT_DIR=./$(HOST_AGENT_DIR) $(MAKE) host-agent-binary

host-agent-binary: $(RELEASE_DIR)
	docker run \
		--rm \
		-e CGO_ENABLED=0 \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-v "$$(pwd):/workspace$(DOCKER_VOL_OPTS)" \
		-w /workspace \
		golang:1.25 \
		go build -buildvcs=false -a -ldflags "$(GOLDFLAGS)" \
		-o ./bin/$(notdir $(RELEASE_BINARY))-$(GOOS)-$(GOARCH) $(HOST_AGENT_DIR)


##@ Release

$(RELEASE_DIR):
	rm -rf $(RELEASE_DIR)
	mkdir -p $(RELEASE_DIR)

build-release-artifacts: build-cluster-templates build-infra-yaml build-metadata-yaml build-host-agent-binary ## Builds release artifacts

build-cluster-templates: $(RELEASE_DIR) cluster-templates
	cp $(BYOH_TEMPLATES)/v1beta1/templates/docker/cluster-template.yaml $(RELEASE_DIR)/cluster-template-docker.yaml
	cp $(BYOH_TEMPLATES)/v1beta1/templates/docker/cluster-template-topology-docker.yaml $(RELEASE_DIR)/cluster-template-topology-docker.yaml
	cp $(BYOH_TEMPLATES)/v1beta1/templates/docker/clusterclass-quickstart-docker.yaml $(RELEASE_DIR)/clusterclass-quickstart-docker.yaml
	cp $(BYOH_TEMPLATES)/v1beta1/templates/vm/cluster-template.yaml $(RELEASE_DIR)/cluster-template.yaml
	cp $(BYOH_TEMPLATES)/v1beta1/templates/vm/cluster-template-topology.yaml $(RELEASE_DIR)/cluster-template-topology.yaml
	cp $(BYOH_TEMPLATES)/v1beta1/templates/vm/clusterclass-quickstart.yaml $(RELEASE_DIR)/clusterclass-quickstart.yaml


build-infra-yaml: ## Generate infrastructure-components.yaml for the provider
	cd config/manager && $(KUSTOMIZE) edit set image ghcr.io/cohesity/cluster-api-provider-bringyourownhost-controller=${IMG}
	$(KUSTOMIZE) build config/default > $(RELEASE_DIR)/infrastructure-components.yaml

build-metadata-yaml:
	cp metadata.yaml $(RELEASE_DIR)/metadata.yaml

build-host-agent-binary: host-agent-binaries
	cp bin/byoh-hostagent-linux-amd64 $(RELEASE_DIR)/byoh-hostagent-linux-amd64

clean:
	git clean -xfd

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUBECTL ?= kubectl
KIND ?= kind
KUSTOMIZE ?= go run sigs.k8s.io/kustomize/kustomize/v5
CONTROLLER_GEN ?= go run sigs.k8s.io/controller-tools/cmd/controller-gen
GINKGO ?= go run $(GINKGO_PKG)
YQ ?= go run github.com/mikefarah/yq/v4
ENVTEST ?= $(LOCALBIN)/setup-envtest
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint

## Tool Versions
KUSTOMIZE_VERSION ?= v5.6.0
CONTROLLER_TOOLS_VERSION ?= v0.18.0
#ENVTEST_VERSION is the version of controller-runtime release branch to fetch the envtest setup script (i.e. release-0.20)
ENVTEST_VERSION ?= $(shell go list -m -f "{{ .Version }}" sigs.k8s.io/controller-runtime | awk -F'[v.]' '{printf "release-%d.%d", $$2, $$3}')
#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries (i.e. 1.31)
ENVTEST_K8S_VERSION ?= $(shell go list -m -f "{{ .Version }}" k8s.io/api | awk -F'[v.]' '{printf "1.%d", $$3}')
GOLANGCI_LINT_VERSION ?= v2.4.0

.PHONY: setup-envtest
setup-envtest: envtest ## Download the binaries required for ENVTEST in the local bin directory.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] && [ "$$(readlink -- "$(1)" 2>/dev/null)" = "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $$(realpath $(1)-$(3)) $(1)
endef
