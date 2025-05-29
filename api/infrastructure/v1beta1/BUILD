load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "v1beta1",
    srcs = [
        "bootstrapkubeconfig_types.go",
        "byocluster_types.go",
        "byoclustertemplate_types.go",
        "byohost_types.go",
        "byomachine_types.go",
        "byomachinetemplate_types.go",
        "condition_consts.go",
        "doc.go",
        "groupversion_info.go",
        "k8sinstallerconfig_types.go",
        "k8sinstallerconfigtemplate_types.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "github.com/cohesity/cluster-api-provider-bringyourownhost/api/infrastructure/v1beta1",
    visibility = ["//visibility:public"],
    deps = [
        "@io_k8s_api//core/v1:core",
        "@io_k8s_apimachinery//pkg/apis/meta/v1:meta",
        "@io_k8s_apimachinery//pkg/runtime",
        "@io_k8s_apimachinery//pkg/runtime/schema",
        "@io_k8s_sigs_cluster_api//api/v1beta1",
        "@io_k8s_sigs_controller_runtime//pkg/scheme",
    ],
)
