load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "reconciler",
    srcs = [
        "doc.go",
        "host_reconciler.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/agent/reconciler",
    visibility = ["//visibility:public"],
    deps = [
        "//agent/cloudinit",
        "//agent/registration",
        "//apis/infrastructure/v1beta1",
        "//common",
        "@com_github_kube_vip_kube_vip//pkg/vip",
        "@com_github_pkg_errors//:errors",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_client_go//tools/record",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//util/conditions",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_cluster_api//util/predicates",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/manager",
    ],
)

go_test(
    name = "reconciler_test",
    srcs = [
        "reconciler_suite_test.go",
        "reconciler_test.go",
    ],
    deps = [
        ":reconciler",
        "//agent/cloudinit/cloudinitfakes",
        "//apis/infrastructure/v1beta1",
        "//test/builder",
        "//test/utils/events",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_client_go//rest",
        "@io_k8s_client_go//tools/record",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//util/conditions",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/envtest",
        "@io_k8s_sigs_controller_runtime//pkg/manager",
    ],
)
