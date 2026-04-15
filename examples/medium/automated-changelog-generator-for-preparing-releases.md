# Requirement: "an automated changelog generator for preparing releases"

Collects change entries from a directory, groups them by kind, and renders a release section. File I/O and time go through thin std primitives.

std
  std.fs
    std.fs.list_dir
      fn (path: string) -> result[list[string], string]
      + returns file names in the directory
      - returns error when path does not exist
      # filesystem
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads the full file contents as a string
      - returns error on missing file or permission denied
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes contents to path, replacing any existing file
      # filesystem
  std.time
    std.time.today_iso
      fn () -> string
      + returns the current date as YYYY-MM-DD
      # time

changelog
  changelog.entry_parse
    fn (raw: string) -> result[change_entry, string]
    + parses a single change entry with kind and body fields
    - returns error when the kind field is missing
    # parsing
  changelog.collect_unreleased
    fn (dir: string) -> result[list[change_entry], string]
    + reads all entry files in dir and returns the parsed entries
    - returns error when any entry fails to parse
    # aggregation
    -> std.fs.list_dir
    -> std.fs.read_all
  changelog.group_by_kind
    fn (entries: list[change_entry]) -> map[string, list[change_entry]]
    + returns entries grouped by their kind (added, changed, fixed, removed)
    + preserves input order within each group
    # grouping
  changelog.render_release
    fn (version: string, groups: map[string, list[change_entry]]) -> string
    + returns a markdown section titled with version and date
    + emits one subsection per non-empty kind
    ? kind order is fixed: added, changed, fixed, removed
    # rendering
    -> std.time.today_iso
  changelog.prepend_release
    fn (changelog_path: string, release_section: string) -> result[void, string]
    + inserts the new release section near the top of an existing changelog file
    - returns error when the file cannot be read or written
    # update
    -> std.fs.read_all
    -> std.fs.write_all
