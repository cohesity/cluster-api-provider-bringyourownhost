load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "feature",
    srcs = ["feature.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/feature",
    visibility = ["//visibility:public"],
    deps = [
        "@io_k8s_apimachinery//pkg/util/runtime",
        "@io_k8s_component_base//featuregate",
    ],
)
