"Orion extension for generating go_image targets."

aspect.register_rule_kind("go_image", {
    "From": "//tools/oci:go_image.bzl",
    "ResolveAttrs": ["binary", "visibility"],
})

def declare_targets(ctx):
    for file in ctx.sources:
        if len(file.query_results["has_main"]) > 0:
            go_library_target_name = path.base(ctx.rel)
            ctx.targets.add(
                name = "image",
                kind = "go_image",
                attrs = {
                    "binary": go_library_target_name,
                    "visibility": attr.string_list(),
                },
            )

aspect.register_configure_extension(
    id = "go_image",
    prepare = lambda cfg: aspect.PrepareResult(
        sources = [aspect.SourceExtensions(".go")],
        queries = {
            # orion does not have a go grammar yet.
            "has_main": aspect.RegexQuery(
                filter = "*.go",
                expression = """(?P<main>func\\s+main\\(.*?\\))""",
            ),
        },
    ),
    declare = declare_targets,
)
