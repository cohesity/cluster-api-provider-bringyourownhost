load("@bazel_skylib//lib:paths.bzl", "paths")
load("@rules_multirun//:defs.bzl", "command", "multirun")
load("@rules_shell//shell:sh_binary.bzl", "sh_binary")

"""Targets in the repository root"""

# We prefer BUILD instead of BUILD.bazel
# gazelle:build_file_name BUILD

load("@gazelle//:def.bzl", "gazelle")

# keep
alias(
    name = "format",
    actual = "//tools/format",
)

exports_files(
    [".shellcheckrc"],
    visibility = ["//:__subpackages__"],
)

# gazelle:prefix github.com/cohesity/cluster-api-provider-bringyourownhost
gazelle(
    name = "gazelle",
    gazelle = "@multitool//tools/gazelle",
)

sh_binary(
    name = "kubebuilder-setup",
    srcs = ["kubebuilder-setup.sh"],
)

alias(
    name = "controller-image",
    actual = "//cmd:image",
)

alias(
    name = "host-agent-binaries",
    actual = "//agent",
)

command(
    name = "generate-crds",
    arguments = [
        "rbac:roleName=manager-role",
        "crd",
        "webhook",
        'paths="./..."',
        "output:crd:artifacts:config=config/crd/bases",
    ],
    command = "@io_k8s_sigs_controller_tools//cmd/controller-gen",
    data = [
        "@rules_go//go",
    ],
    run_from_workspace_root = True,
)

command(
    name = "generate-code",
    arguments = [
        'object:headerFile="hack/boilerplate.go.txt"',
        'paths="./..."',
    ],
    command = "@io_k8s_sigs_controller_tools//cmd/controller-gen",
    data = [
        "@rules_go//go",
    ],
    run_from_workspace_root = True,
)

crd_files = glob(
    ["config/crd/bases/*.yaml"],
)

[command(
    name = "yq-fix-" + paths.basename(crd_file),
    arguments = [
        "-i",
        "eval",
        "del(.metadata.creationTimestamp)",
        crd_file,
    ],
    command = "@com_github_mikefarah_yq_v4//:v4",
    run_from_workspace_root = True,
) for crd_file in crd_files]

multirun(
    name = "yq-fix",
    commands = [":yq-fix-{}".format(paths.basename(crd_file)) for crd_file in crd_files],
    jobs = 0,  # Set to 0 to run in parallel, defaults to sequential
    visibility = ["//visibility:public"],
)

multirun(
    name = "generate",
    commands = [
        ":generate-crds",
        ":yq-fix",
        ":generate-code",
    ],
    jobs = 0,  # Set to 0 to run in parallel, defaults to sequential
    visibility = ["//visibility:public"],
)
