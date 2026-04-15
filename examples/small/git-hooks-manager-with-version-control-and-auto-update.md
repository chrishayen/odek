# Requirement: "a library for managing per-repo and shared git hooks with version control and auto update"

Installs, lists, and updates git hook scripts maintained in a shared source, mirrored into each repo's hooks directory.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns the entries directly under path
      - returns error when path does not exist
      # filesystem
    std.fs.copy_file
      fn (src: string, dst: string) -> result[void, string]
      + copies src to dst, creating parent directories as needed
      - returns error when src does not exist
      # filesystem

git_hooks
  git_hooks.discover
    fn (repo_root: string, shared_root: string) -> result[hook_set, string]
    + returns the union of hooks available in the repo's own hooks directory and the shared source
    - returns error when repo_root is not a git repository
    # discovery
    -> std.fs.list_dir
  git_hooks.install
    fn (repo_root: string, set: hook_set) -> result[void, string]
    + writes each hook into .git/hooks with the executable bit set
    - returns error when .git/hooks is not writable
    # installation
    -> std.fs.copy_file
  git_hooks.update
    fn (repo_root: string, shared_root: string) -> result[list[string], string]
    + refreshes installed hooks from the shared source and returns the names that changed
    ? hooks whose content is identical are left untouched
    # update
    -> std.fs.copy_file
  git_hooks.is_up_to_date
    fn (repo_root: string, shared_root: string) -> result[bool, string]
    + returns true when every installed hook matches the shared source
    # status
