# Requirement: "a library for managing multiple source repositories"

Declare a set of repositories in a config, then clone, sync, and run commands across them.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns file contents
      - returns error when unreadable
      # filesystem
    std.fs.exists
      fn (path: string) -> bool
      + returns true when the path exists
      # filesystem
  std.process
    std.process.run
      fn (cwd: string, program: string, args: list[string]) -> result[string, string]
      + runs the program and returns stdout on zero exit
      - returns error including stderr on non-zero exit
      # process

multirepo
  multirepo.parse_config
    fn (raw: string) -> result[list[repo_spec], string]
    + parses a config listing repo name, url, and local path
    - returns error on missing required fields
    # config
  multirepo.clone_missing
    fn (specs: list[repo_spec]) -> result[list[string], string]
    + clones any repo whose local path does not exist and returns the cloned names
    - returns error when a clone fails
    # clone
    -> std.fs.exists
    -> std.process.run
  multirepo.sync_all
    fn (specs: list[repo_spec]) -> result[map[string,string], string]
    + pulls each repo and returns a name-to-status map
    - returns error when any repo has local changes preventing pull
    # sync
    -> std.process.run
  multirepo.run_in_each
    fn (specs: list[repo_spec], program: string, args: list[string]) -> map[string, result[string,string]]
    + runs the command in each repo and returns a per-repo result
    + continues past failures so every repo reports
    # fanout
    -> std.process.run
  multirepo.filter_by_tag
    fn (specs: list[repo_spec], tag: string) -> list[repo_spec]
    + returns specs whose tag list contains the given tag
    + empty list when nothing matches
    # filtering
