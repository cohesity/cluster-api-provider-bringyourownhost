load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "bootstraptoken",
    srcs = ["token.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/common/bootstraptoken",
    visibility = ["//visibility:public"],
    deps = [
        "//apis/infrastructure/v1beta1",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_client_go//tools/clientcmd/api",
        "@io_k8s_cluster_bootstrap//token/api",
        "@io_k8s_cluster_bootstrap//token/util",
    ],
)
