load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "v1beta1",
    srcs = [
        "bootstrapkubeconfig_webhook.go",
        "byohost_webhook.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/internal/webhook/infrastructure/v1beta1",
    visibility = ["//:__subpackages__"],
    deps = [
        "//api/infrastructure/v1beta1",
        "@io_k8s_api//admission/v1:admission",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_apimachinery//pkg/util/validation/field",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/webhook",
        "@io_k8s_sigs_controller_runtime//pkg/webhook/admission",
    ],
)

go_test(
    name = "v1beta1_test",
    srcs = [
        "bootstrapkubeconfig_webhook_test.go",
        "byohost_webhook_internal_test.go",
        "byohost_webhook_test.go",
        "webhook_suite_test.go",
    ],
    deps = [
        ":v1beta1",
        "//api/infrastructure/v1beta1",
        "//test/builder",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
        "@io_k8s_api//admission/v1:admission",
        "@io_k8s_api//authentication/v1:authentication",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_client_go//rest",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/envtest",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
        "@io_k8s_sigs_controller_runtime//pkg/metrics/server",
        "@io_k8s_sigs_controller_runtime//pkg/webhook",
        "@io_k8s_sigs_controller_runtime//pkg/webhook/admission",
    ],
)
