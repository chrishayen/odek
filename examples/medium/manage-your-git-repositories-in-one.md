# Requirement: "a library to manage multiple version-control repositories in one place"

Tracks a set of local repository paths and performs batch operations (status, fetch, pull) over them via a pluggable VCS backend.

std
  std.fs
    std.fs.exists
      @ (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.process
    std.process.run_in
      @ (cwd: string, command: string, args: list[string]) -> result[tuple[i32, string], string]
      + runs a command in a working directory and returns exit code and stdout
      - returns error when the command cannot be launched
      # process

repo_manager
  repo_manager.new
    @ () -> manager_state
    + creates an empty repository set
    # construction
  repo_manager.add
    @ (state: manager_state, path: string) -> result[manager_state, string]
    + registers a local repository path
    - returns error when the path does not exist
    - returns error when the path is already registered
    # registration
    -> std.fs.exists
  repo_manager.remove
    @ (state: manager_state, path: string) -> result[manager_state, string]
    + unregisters a repository path
    - returns error when the path is not registered
    # registration
  repo_manager.status_one
    @ (state: manager_state, path: string) -> result[repo_status, string]
    + returns current branch, ahead/behind counts, and dirty flag for one repository
    - returns error when the path is not a repository
    # status
    -> std.process.run_in
  repo_manager.status_all
    @ (state: manager_state) -> list[tuple[string, repo_status]]
    + returns status for every registered repository
    + repositories that fail appear with an error-marked status
    # status
  repo_manager.fetch_all
    @ (state: manager_state) -> list[tuple[string, result[void, string]]]
    + fetches each repository and returns per-path results
    # operations
    -> std.process.run_in
  repo_manager.pull_all
    @ (state: manager_state) -> list[tuple[string, result[void, string]]]
    + pulls each repository and returns per-path results
    # operations
    -> std.process.run_in
