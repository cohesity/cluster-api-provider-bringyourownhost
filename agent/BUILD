load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("//tools/oci:go_image.bzl", "go_image")

go_library(
    name = "agent_lib",
    srcs = ["main.go"],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/agent",
    visibility = ["//visibility:private"],
    deps = [
        "//agent/cloudinit",
        "//agent/reconciler",
        "//agent/registration",
        "//agent/version",
        "//api/infrastructure/v1beta1",
        "//feature",
        "@com_github_go_logr_logr//:logr",
        "@com_github_spf13_pflag//:pflag",
        "@io_k8s_api//certificates/v1:certificates",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/fields",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_client_go//rest",
        "@io_k8s_klog_v2//:klog",
        "@io_k8s_klog_v2//klogr",
        "@io_k8s_sigs_cluster_api//api/core/v1beta1",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/cache",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/metrics/server",
        "@io_k8s_sigs_controller_runtime//pkg/webhook",
    ],
)

go_binary(
    name = "agent",
    embed = [":agent_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "agent_test",
    srcs = [
        "help_flag_test.go",
        "host_agent_suite_test.go",
        "host_agent_test.go",
        "label_flag_test.go",
    ],
    embed = [":agent_lib"],
    deps = [
        "//agent/registration",
        "//agent/version",
        "//api/infrastructure/v1beta1",
        "//test/builder",
        "//test/e2e",
        "@com_github_docker_docker//api/types",
        "@com_github_docker_docker//api/types/container",
        "@com_github_docker_docker//client",
        "@com_github_jackpal_gateway//:gateway",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
        "@com_github_onsi_gomega//gexec",
        "@io_k8s_api//certificates/v1:certificates",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//rest",
        "@io_k8s_klog_v2//klogr",
        "@io_k8s_sigs_cluster_api//api/core/v1beta1",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/envtest",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
    ],
)

go_image(
    name = "image",
    binary = "agent",
    visibility = ["//visibility:public"],
)
