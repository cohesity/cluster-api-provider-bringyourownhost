load("@rules_shell//shell:sh_binary.bzl", "sh_binary")

"""Targets in the repository root"""

# We prefer BUILD instead of BUILD.bazel
# gazelle:build_file_name BUILD

load("@gazelle//:def.bzl", "gazelle")

alias(
    name = "format",
    actual = "//tools/format",
)

exports_files(
    [".shellcheckrc"],
    visibility = ["//:__subpackages__"],
)

# gazelle:prefix github.com/cohesity/cluster-api-provider-bringyourownhost
gazelle(name = "gazelle")

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
