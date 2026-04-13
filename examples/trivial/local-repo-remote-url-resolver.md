# Requirement: "resolve the remote hosting URL of a repository in a local directory"

Reads the origin remote of a local repository and returns the equivalent web URL. Opening it is the caller's responsibility.

std
  std.process
    std.process.run
      @ (cmd: string, args: list[string]) -> result[process_output, string]
      + runs a command and returns stdout, stderr, and exit code
      - returns error when the binary cannot be located
      # process

repo_url
  repo_url.origin_web_url
    @ (directory: string) -> result[string, string]
    + returns an https web URL derived from the origin remote
    + converts ssh-style git@host:owner/repo into https://host/owner/repo
    - returns error when the directory is not a repository
    - returns error when there is no origin remote
    # url_resolution
    -> std.process.run
