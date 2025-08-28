"""Definition of linter tools

This is in its own package because it has so many loading-time symbols,
we don't want to trigger eager fetches of these for builds which aren't running linters.
"""

load("@rules_go//go:def.bzl", "TOOLS_NOGO", "nogo")

package(default_visibility = ["//:__subpackages__"])

nogo(
    name = "nogo",
    config = "nogo_config.json",
    visibility = ["//visibility:public"],
    deps = TOOLS_NOGO + [
        # Additional analyzers not included in TOOLS_NOGO
        "@org_golang_x_tools//go/analysis/passes/fieldalignment",
        "@org_golang_x_tools//go/analysis/passes/gofix",
        "@org_golang_x_tools//go/analysis/passes/hostport",
        "@org_golang_x_tools//go/analysis/passes/httpmux",
        "@org_golang_x_tools//go/analysis/passes/pkgfact",
        "@org_golang_x_tools//go/analysis/passes/reflectvaluecompare",
        "@org_golang_x_tools//go/analysis/passes/shadow",
        "@org_golang_x_tools//go/analysis/passes/stdversion",
    ],
)
