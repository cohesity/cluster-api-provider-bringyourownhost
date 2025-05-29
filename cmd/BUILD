load("@rules_go//go:def.bzl", "go_binary", "go_library")
load("//tools/oci:go_image.bzl", "go_image")

go_library(
    name = "cmd_lib",
    srcs = ["main.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/cmd",
    visibility = ["//visibility:private"],
    deps = [
        "//api/infrastructure/v1beta1",
        "//internal/controller/infrastructure",
        "//internal/webhook/infrastructure/v1beta1",
        "@io_k8s_api//admission/v1:admission",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/util/runtime",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_client_go//plugin/pkg/client/auth",
        "@io_k8s_klog_v2//:klog",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//controllers/remote",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/certwatcher",
        "@io_k8s_sigs_controller_runtime//pkg/controller",
        "@io_k8s_sigs_controller_runtime//pkg/healthz",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
        "@io_k8s_sigs_controller_runtime//pkg/metrics/filters",
        "@io_k8s_sigs_controller_runtime//pkg/metrics/server",
        "@io_k8s_sigs_controller_runtime//pkg/webhook",
    ],
)

go_binary(
    name = "cmd",
    embed = [":cmd_lib"],
    visibility = ["//visibility:public"],
)

go_image(
    name = "image",
    binary = "cmd",
    visibility = ["//visibility:public"],
)
