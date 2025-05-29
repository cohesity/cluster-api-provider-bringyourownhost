load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "cloudinitfakes",
    srcs = [
        "fake_icmd_runner.go",
        "fake_ifile_writer.go",
        "fake_itemplate_parser.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/agent/cloudinit/cloudinitfakes",
    visibility = ["//visibility:public"],
    deps = ["//agent/cloudinit"],
)
