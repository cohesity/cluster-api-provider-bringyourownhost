load("@rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_shell//shell:sh_binary.bzl", "sh_binary")

"""Targets in the repository root"""

# We prefer BUILD instead of BUILD.bazel
# gazelle:build_file_name BUILD

load("@gazelle//:def.bzl", "gazelle")

alias(
    name = "format",
    actual = "//tools/format",
)

exports_files(
    [".shellcheckrc"],
    visibility = ["//:__subpackages__"],
)

# gazelle:prefix github.com/cohesity/cluster-api-provider-bringyourownhost
gazelle(name = "gazelle")

sh_binary(
    name = "kubebuilder-setup",
    srcs = ["kubebuilder-setup.sh"],
)

alias(
    name = "controller-image",
    actual = "//cmd:image",
)

alias(
    name = "host-agent-binaries",
    actual = "//agent",
)

go_library(
    name = "cluster-api-provider-bringyourownhost_lib",
    srcs = ["main.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost",
    visibility = ["//visibility:private"],
    deps = [
        "//apis/infrastructure/v1beta1",
        "//controllers/infrastructure",
        "@io_k8s_api//admission/v1beta1",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/util/runtime",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_client_go//plugin/pkg/client/auth",
        "@io_k8s_klog_v2//:klog",
        "@io_k8s_klog_v2//klogr",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//controllers/remote",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/controller",
        "@io_k8s_sigs_controller_runtime//pkg/healthz",
        "@io_k8s_sigs_controller_runtime//pkg/webhook",
    ],
)

go_binary(
    name = "cluster-api-provider-bringyourownhost",
    embed = [":cluster-api-provider-bringyourownhost_lib"],
    visibility = ["//visibility:public"],
)
