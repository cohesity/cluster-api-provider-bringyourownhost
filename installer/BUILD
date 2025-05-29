load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "installer",
    srcs = [
        "bundle_downloader.go",
        "installer.go",
        "registry.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/installer",
    visibility = ["//visibility:public"],
    deps = [
        "//installer/internal/algo",
        "@com_github_go_logr_logr//:logr",
    ],
)

go_test(
    name = "installer_test",
    srcs = [
        "bundle_downloader_test.go",
        "installer_suite_test.go",
        "installer_test.go",
        "registry_internal_test.go",
    ],
    embed = [":installer"],
    deps = [
        "@com_github_go_logr_logr//:logr",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
    ],
)
