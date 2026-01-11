load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "runtime",
    srcs = [
        "impl.go",
        "runtime.go",
        "runtime_unix.go",
        "runtime_windows.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/util/runtime",
    visibility = ["//visibility:public"],
    deps = [
        "@io_k8s_apimachinery//pkg/util/errors",
        "@io_k8s_cri_api//pkg/apis",
        "@io_k8s_cri_api//pkg/apis/runtime/v1:runtime",
        "@io_k8s_cri_client//pkg",
        "@io_k8s_klog_v2//:klog",
    ] + select({
        "@rules_go//go/platform:windows": [
            "@com_github_microsoft_go_winio//:go-winio",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "runtime_test",
    srcs = [
        "runtime_fake_test.go",
        "runtime_test.go",
    ],
    embed = [":runtime"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@io_k8s_cri_api//pkg/apis",
        "@io_k8s_cri_api//pkg/apis/runtime/v1:runtime",
    ],
)
