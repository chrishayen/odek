# Requirement: "a patch-based distributed version control system"

A version control model where history is an unordered set of patches that commute when they do not overlap. The project layer manages patches, application, and merging.

std
  std.fs
    std.fs.read_bytes
      fn (path: string) -> result[bytes, string]
      + returns the full contents of the file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_bytes
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to the file
      - returns error when the path is not writable
      # filesystem
  std.hash
    std.hash.blake3
      fn (data: bytes) -> bytes
      + returns a 32-byte BLAKE3 hash of data
      # hashing
  std.time
    std.time.now_seconds
      fn () -> i64
      + returns current unix time in seconds
      # time

pvcs
  pvcs.init_repo
    fn (root: string) -> result[repo_state, string]
    + creates an empty repository rooted at root
    - returns error when root already contains a repository
    # construction
  pvcs.diff_text
    fn (before: string, after: string) -> list[hunk]
    + returns a list of line-level hunks representing the change
    ? identical inputs produce an empty list
    # diffing
  pvcs.make_patch
    fn (author: string, message: string, changes: list[file_change]) -> patch
    + assembles a patch with a stable content-addressed id
    # patch_creation
    -> std.hash.blake3
    -> std.time.now_seconds
  pvcs.apply_patch
    fn (state: repo_state, p: patch) -> result[repo_state, string]
    + applies a patch whose dependencies are already present in the repo
    - returns error when a dependency is missing
    - returns error when a text hunk does not match the current content
    # application
  pvcs.patches_commute
    fn (a: patch, b: patch) -> bool
    + returns true when the two patches touch disjoint regions of disjoint files
    ? commuting patches can be applied in either order with the same result
    # commutation
  pvcs.unapply_patch
    fn (state: repo_state, patch_id: bytes) -> result[repo_state, string]
    + reverses a patch, producing the state that would exist without it
    - returns error when the patch has dependents still applied
    # inversion
  pvcs.merge
    fn (local: repo_state, remote: repo_state) -> result[repo_state, list[conflict]]
    + returns a merged repo containing the union of both patch sets
    - returns the list of conflicts when overlapping patches disagree
    # merging
  pvcs.checkout_file
    fn (state: repo_state, path: string) -> result[string, string]
    + returns the text content of a file as currently materialized
    - returns error when the path is not tracked
    # checkout
    -> std.fs.read_bytes
  pvcs.record_working_tree
    fn (state: repo_state, author: string, message: string) -> result[patch, string]
    + scans the working tree for changes and produces a patch capturing them
    - returns error when no changes are present
    # recording
    -> std.fs.write_bytes
  pvcs.missing_dependencies
    fn (state: repo_state, p: patch) -> list[bytes]
    + returns the ids of patches p depends on that are not yet applied
    # dependency_check
  pvcs.patch_id
    fn (p: patch) -> bytes
    + returns the content-addressed id of a patch
    # identity
    -> std.hash.blake3
