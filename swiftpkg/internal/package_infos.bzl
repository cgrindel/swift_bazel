"""API for creating and loading Swift package information."""

load(":repository_utils.bzl", "repository_utils")

def _get_dump_manifest(repository_ctx, env = {}, working_directory = ""):
    """Returns a dict representing the package dump for an SPM package.

    Args:
        repository_ctx: A `repository_ctx`.
        env: A `dict` of environment variables that will be included in the
             command execution.
        working_directory: A `string` specifying the directory for the SPM package.

    Returns:
        A `dict` representing an SPM package dump.
    """
    return repository_utils.parsed_json_from_spm_command(
        repository_ctx,
        ["swift", "package", "dump-package"],
        env = env,
        working_directory = working_directory,
        debug_json_path = "dump.json",
    )

def _get_desc_manifest(repository_ctx, env = {}, working_directory = ""):
    """Returns a dict representing the package description for an SPM package.

    Args:
        repository_ctx: A `repository_ctx`.
        env: A `dict` of environment variables that will be included in the
             command execution.
        working_directory: A `string` specifying the directory for the SPM package.

    Returns:
        A `dict` representing an SPM package description.
    """
    return repository_utils.parsed_json_from_spm_command(
        repository_ctx,
        ["swift", "package", "describe", "--type", "json"],
        env = env,
        working_directory = working_directory,
        debug_json_path = "desc.json",
    )

def _get(repository_ctx, directory, env = {}):
    dump_manifest = _get_dump_manifest(
        repository_ctx,
        env = env,
        working_directory = directory,
    )
    desc_manifest = _get_desc_manifest(
        repository_ctx,
        env = env,
        working_directory = directory,
    )
    return _new(
        dump_manifest = dump_manifest,
        desc_manifest = desc_manifest,
    )

def _new_dependency_requirement_from_desc_json_map(reqMap):
    ranges = reqMap.get("range")
    if ranges != None:
        return _new_dependency_requirement(
            ranges = [
                _new_version_range(
                    lower = rangeMap["lower_bound"],
                    upper = rangeMap["upper_bound"],
                )
                for rangeMap in ranges
            ],
        )
    return None

def _new_dependency_from_desc_json_map(depMap):
    return _new_dependency(
        identity = depMap["identity"],
        type = depMap["type"],
        url = depMap["url"],
        requirement = _new_dependency_requirement_from_desc_json_map(depMap["requirement"]),
    )

def _new_from_parsed_json(dump_manifest, desc_manifest):
    tools_version = dump_manifest["toolsVersion"]["_version"]
    platforms = [
        _new_platform(name = pl["platformName"], version = pl["version"])
        for pl in dump_manifest["platforms"]
    ]
    dependencies = [
        _new_dependency_from_desc_json_map(depMap)
        for depMap in desc_manifest["dependencies"]
    ]

    return _new(
        name = dump_manifest["name"],
        path = desc_manifest["path"],
        tools_version = tools_version,
        platforms = platforms,
        dependencies = dependencies,
    )

# def _new(directory, dump_manifest, desc_manifest):
#     return struct(
#         directory = directory,
#         dump_manifest = dump_manifest,
#         desc_manifest = desc_manifest,
#     )

def _new(
        name,
        path,
        tools_version = None,
        platforms = [],
        dependencies = [],
        products = [],
        targets = []):
    return struct(
        name = name,
        path = path,
        tools_version = tools_version,
        platforms = platforms,
        dependencies = dependencies,
        products = products,
        targets = targets,
    )

def _new_platform(name, version):
    return struct(
        name = name,
        version = version,
    )

def _new_dependency(identity, type, url, requirement):
    return struct(
        identity = identity,
        type = type,
        url = url,
        requirement = requirement,
    )

def _new_dependency_requirement(ranges = None):
    return struct(
        ranges = ranges,
    )

def _new_version_range(lower, upper):
    return struct(
        lower = lower,
        upper = upper,
    )

# def _new_product(name, type, targets):
#     return struct(
#         name = name,
#         type = type,
#         targets = targets,
#     )

# def _target(name, type, dependencies = []):
#     return struct(
#         name = name,
#         type = type,
#         dependencies = dependencies,
#     )

package_infos = struct(
    get = _get,
    new = _new,
    new_from_parsed_json = _new_from_parsed_json,
    new_platform = _new_platform,
    new_dependency = _new_dependency,
    new_dependency_requirement = _new_dependency_requirement,
    new_version_range = _new_version_range,
)
