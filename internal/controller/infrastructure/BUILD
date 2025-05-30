load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "infrastructure",
    srcs = [
        "bootstrapkubeconfig_controller.go",
        "byoadmission_controller.go",
        "byocluster_controller.go",
        "byohost_controller.go",
        "byomachine_controller.go",
        "byomachine_scope.go",
        "byomachinetemplate_controller.go",
        "doc.go",
        "k8sinstallerconfig_controller.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/internal/controller/infrastructure",
    visibility = ["//:__subpackages__"],
    deps = [
        "//api/infrastructure/v1beta1",
        "//common/bootstraptoken",
        "//installer",
        "@com_github_go_logr_logr//:logr",
        "@com_github_pkg_errors//:errors",
        "@io_k8s_api//certificates/v1:certificates",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/apis/meta/v1/unstructured",
        "@io_k8s_apimachinery//pkg/labels",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_apimachinery//pkg/selection",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//tools/clientcmd/api/latest",
        "@io_k8s_client_go//tools/record",
        "@io_k8s_cluster_bootstrap//token/util",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//controllers/external",
        "@io_k8s_sigs_cluster_api//controllers/remote",
        "@io_k8s_sigs_cluster_api//util",
        "@io_k8s_sigs_cluster_api//util/annotations",
        "@io_k8s_sigs_cluster_api//util/conditions",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_cluster_api//util/predicates",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/builder",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/controller/controllerutil",
        "@io_k8s_sigs_controller_runtime//pkg/event",
        "@io_k8s_sigs_controller_runtime//pkg/handler",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/predicate",
        "@io_k8s_sigs_controller_runtime//pkg/reconcile",
        "@io_k8s_utils//pointer",
    ],
)

go_test(
    name = "infrastructure_test",
    srcs = [
        "bootstrapkubeconfig_controller_test.go",
        "byoadmission_controller_test.go",
        "byocluster_controller_test.go",
        "byohost_controller_test.go",
        "byomachine_controller_test.go",
        "byomachinetemplate_controller_test.go",
        "k8sinstallerconfig_controller_test.go",
        "suite_test.go",
    ],
    deps = [
        ":infrastructure",
        "//api/infrastructure/v1beta1",
        "//test/builder",
        "//test/utils/events",
        "@com_github_go_logr_logr//:logr",
        "@com_github_onsi_ginkgo_v2//:ginkgo",
        "@com_github_onsi_gomega//:gomega",
        "@io_k8s_api//certificates/v1:certificates",
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/api/errors",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/types",
        "@io_k8s_client_go//kubernetes/fake",
        "@io_k8s_client_go//kubernetes/scheme",
        "@io_k8s_client_go//rest",
        "@io_k8s_client_go//tools/clientcmd",
        "@io_k8s_client_go//tools/record",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_cluster_api//bootstrap/kubeadm/api/v1beta1",
        "@io_k8s_sigs_cluster_api//controllers/remote",
        "@io_k8s_sigs_cluster_api//util",
        "@io_k8s_sigs_cluster_api//util/annotations",
        "@io_k8s_sigs_cluster_api//util/conditions",
        "@io_k8s_sigs_cluster_api//util/patch",
        "@io_k8s_sigs_controller_runtime//:controller-runtime",
        "@io_k8s_sigs_controller_runtime//pkg/client",
        "@io_k8s_sigs_controller_runtime//pkg/client/fake",
        "@io_k8s_sigs_controller_runtime//pkg/controller/controllerutil",
        "@io_k8s_sigs_controller_runtime//pkg/envtest",
        "@io_k8s_sigs_controller_runtime//pkg/log",
        "@io_k8s_sigs_controller_runtime//pkg/log/zap",
        "@io_k8s_sigs_controller_runtime//pkg/metrics/server",
        "@io_k8s_sigs_controller_runtime//pkg/reconcile",
        "@io_k8s_utils//pointer",
    ],
)
