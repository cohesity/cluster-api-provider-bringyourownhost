load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "builder",
    srcs = ["builders.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/test/builder",
    visibility = ["//visibility:public"],
    deps = [
        "//api/infrastructure/v1beta1",
        "@io_k8s_api//certificates/v1:certificates",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
    ],
)
