"""Module for generating data from target dependencies created by `package_infos`."""

load("@cgrindel_bazel_starlib//bzllib:defs.bzl", "bazel_labels")
load(":bazel_repo_names.bzl", "bazel_repo_names")
load(":pkginfo_ext_deps.bzl", "pkginfo_ext_deps")
load(":pkginfo_targets.bzl", "pkginfo_targets")

def make_pkginfo_target_deps(bazel_labels = bazel_labels):
    def _bazel_label(pkg_info, target_dep):
        if target_dep.by_name:
            # GH009: Need to handle the byName references to external modules. Ugh.
            target = pkginfo_targets.get(pkg_info.targets, target_dep.by_name.target_name)
            label = bazel_labels.new(
                package = target.path,
                name = target_dep.by_name.target_name,
            )
        elif target_dep.product:
            prod_ref = target_dep.product
            ext_dep = pkginfo_ext_deps.find_by_identity(
                pkg_info.dependencies,
                prod_ref.dep_identity,
            )
            repo_name = bazel_repo_names.from_url(ext_dep.url)
            label = bazel_labels.new(
                repository_name = repo_name,
                package = "",
                name = prod_ref.product_name,
            )
        else:
            fail("""\
Unrecognized target dependency while generating a Bazel dependency label.\
""")
        return bazel_labels.normalize(label)

    return struct(
        bazel_label = _bazel_label,
    )

pkginfo_target_deps = make_pkginfo_target_deps(bazel_labels = bazel_labels)
