load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "cloudinit",
    srcs = [
        "cloudinit.go",
        "cmd_runner.go",
        "doc.go",
        "file_writer.go",
        "template_parser.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/agent/cloudinit",
    visibility = ["//visibility:public"],
    deps = [
        "//common",
        "@com_github_pkg_errors//:errors",
        "@io_k8s_sigs_yaml//:yaml",
    ],
)

go_test(
    name = "cloudinit_test",
    srcs = [
        "cloudinit_integration_test.go",
        "cloudinit_suite_test.go",
        "cloudinit_test.go",
        "file_writer_test.go",
    ],
    deps = [
        ":cloudinit",
        "//agent/cloudinit/cloudinitfakes",
        "//agent/registration",
        "//common",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
    ],
)
