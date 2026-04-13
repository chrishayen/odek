# Requirement: "a library for reading and writing keep-a-changelog formatted changelogs"

Parse a changelog into typed releases, add entries to the unreleased section, and render it back.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the whole file into memory
      - returns error when the path does not exist
      # io

changelog
  changelog.parse
    @ (source: string) -> result[changelog_doc, string]
    + parses a changelog into its releases and per-release change groups
    - returns error when the top-level title is missing
    - returns error when a release heading is malformed
    # parsing
  changelog.load
    @ (path: string) -> result[changelog_doc, string]
    + reads and parses a changelog file
    - returns error when the file cannot be read
    # loading
    -> std.fs.read_all
    -> changelog.parse
  changelog.add_entry
    @ (doc: changelog_doc, group: string, entry: string) -> changelog_doc
    + adds an entry to the named group under the unreleased section, creating it if missing
    + valid groups are added, changed, deprecated, removed, fixed, security
    # editing
  changelog.release
    @ (doc: changelog_doc, version: string, date: string) -> result[changelog_doc, string]
    + promotes the unreleased section to a new release with the given version and date
    - returns error when unreleased has no entries
    # editing
  changelog.render
    @ (doc: changelog_doc) -> string
    + renders a changelog back to keep-a-changelog markdown
    # serialization
