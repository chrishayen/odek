# Requirement: "a build orchestration library that compiles a project for multiple target platforms"

Resolves a project manifest, downloads the requested toolchains, and runs a platform-specific compile pipeline producing one artifact per target.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + returns file contents
      - returns error when missing
      # filesystem
    std.fs.write_all
      @ (path: string, data: bytes) -> result[void, string]
      + writes data to path
      # filesystem
    std.fs.list_dir
      @ (path: string) -> result[list[string], string]
      + returns the immediate entries in a directory
      # filesystem
  std.http
    std.http.get
      @ (url: string) -> result[bytes, string]
      + returns the response body
      - returns error on non-200
      # networking
  std.crypto
    std.crypto.sha256_hex
      @ (data: bytes) -> string
      + returns the lowercase hex sha-256 digest
      # cryptography
  std.process
    std.process.run
      @ (argv: list[string], env: map[string,string], cwd: string) -> result[process_result, string]
      + runs a child process and captures stdout, stderr, and exit code
      # process_control

multi_platform_builder
  multi_platform_builder.parse_manifest
    @ (raw: string) -> result[project_manifest, string]
    + parses a project manifest describing sources, targets, and dependencies
    - returns error on malformed input
    # manifest
  multi_platform_builder.resolve_targets
    @ (manifest: project_manifest, requested: list[string]) -> result[list[target_spec], string]
    + returns the target specs for the requested platform names
    - returns error when any name is unknown
    # manifest
  multi_platform_builder.fetch_toolchain
    @ (target: target_spec, cache_dir: string) -> result[string, string]
    + downloads, verifies, and unpacks the toolchain if not already cached
    - returns error when the checksum does not match
    # toolchain
    -> std.http.get
    -> std.crypto.sha256_hex
    -> std.fs.write_all
  multi_platform_builder.discover_sources
    @ (src_dir: string, include_globs: list[string]) -> result[list[string], string]
    + returns the source files matched under src_dir
    # sources
    -> std.fs.list_dir
  multi_platform_builder.plan_build
    @ (manifest: project_manifest, target: target_spec, sources: list[string]) -> build_plan
    + produces an ordered list of compile and link steps for the target
    # planning
  multi_platform_builder.run_step
    @ (step: build_step, toolchain_dir: string, workspace: string) -> result[void, string]
    + executes a single compile or link step and fails on non-zero exit
    - returns error with captured stderr
    # execution
    -> std.process.run
  multi_platform_builder.execute_plan
    @ (plan: build_plan, toolchain_dir: string, workspace: string) -> result[string, string]
    + runs every step in order and returns the path to the produced artifact
    - returns error at the first failing step
    # execution
  multi_platform_builder.build_all
    @ (manifest_path: string, targets: list[string], workspace: string, cache_dir: string) -> result[map[string,string], string]
    + convenience entry point returning a map from target name to artifact path
    # orchestration
    -> std.fs.read_all
  multi_platform_builder.clean
    @ (workspace: string) -> result[void, string]
    + removes all intermediate build output from workspace
    # maintenance
    -> std.fs.list_dir
  multi_platform_builder.incremental_key
    @ (step: build_step, input_hashes: list[string]) -> string
    + returns a cache key combining the step description with its inputs
    # caching
    -> std.crypto.sha256_hex
