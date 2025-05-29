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

# gazelle:prefix github.com/example/project
gazelle(name = "gazelle")
