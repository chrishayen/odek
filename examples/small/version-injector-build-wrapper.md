# Requirement: "a build wrapper that injects version information into compiled artifacts"

Derives version metadata from source-control state and formats it for injection into a build command.

std
  std.process
    std.process.run_capture
      fn (command: string, args: list[string]) -> result[string, string]
      + runs a command and returns its stdout on success
      - returns error containing stderr on non-zero exit
      # process

version_injector
  version_injector.collect
    fn (repo_path: string) -> result[version_info, string]
    + reads commit hash, branch, and tag from the repository at repo_path
    - returns error when repo_path is not a source-control checkout
    # collection
    -> std.process.run_capture
  version_injector.format_flags
    fn (info: version_info, symbol_prefix: string) -> list[string]
    + returns linker flag strings that set symbols under symbol_prefix to the collected values
    + produces one flag per version field
    # formatting
