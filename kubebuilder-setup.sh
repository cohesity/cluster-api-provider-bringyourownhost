#!/usr/bin/env bash
#
# Copyright 2025 Cohesity, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

set -e
set -o errexit
set -o pipefail

current_script=$(readlink -f "$0")
current_dir=$(dirname "${current_script}")

export GO111MODULE=on

KUBEBUILDER=operator-sdk
# KUBEBUILDER=kubebuilder

echo "Initializing ..."

${KUBEBUILDER} init --domain k8s.cohesity.com --repo cohesity/athena/infra-operator
${KUBEBUILDER} init --domain cluster.x-k8s.io --repo github.com/cohesity/cluster-api-provider-bringyourownhost --project-name byoh
${KUBEBUILDER} edit --multigroup=true

echo "Creating hack/boilerplate.go.txt ..."

cat >"hack/boilerplate.go.txt" <<EOF
// Copyright 2025 Cohesity, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
EOF

${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind ByoMachine
${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind ByoHost
${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind ByoCluster
${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind ByoMachineTemplate
${KUBEBUILDER} create api --resource --controller=false --group infrastructure --version v1beta1 --kind ByoClusterTemplate
${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind K8sInstallerConfig
${KUBEBUILDER} create api --resource --controller=false --group infrastructure --version v1beta1 --kind K8sInstallerConfigTemplate
${KUBEBUILDER} create api --resource --controller --group infrastructure --version v1beta1 --kind BootstrapKubeconfig

${KUBEBUILDER} create webhook --group infrastructure --version v1beta1 --kind ByoHost --programmatic-validation
# ${KUBEBUILDER} create webhook --group infrastructure --version v1beta1 --kind ByoCluster --defaulting --programmatic-validation
${KUBEBUILDER} create webhook --group infrastructure --version v1beta1 --kind BootstrapKubeconfig --programmatic-validation
