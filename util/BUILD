load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "util",
    srcs = ["util.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/util",
    visibility = ["//visibility:public"],
    deps = [
        "//agent/cloudinit",
        "//api/infrastructure/v1beta1",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
    ],
)
