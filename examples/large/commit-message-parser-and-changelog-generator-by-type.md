# Requirement: "a framework for parsing structured commit messages and generating a changelog grouped by change type"

Parses conventional-style commit headers, classifies them, groups them by section, and renders a changelog from a template.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits on newline characters
      # strings
    std.strings.trim
      fn (s: string) -> string
      + removes leading and trailing whitespace
      # strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
    std.strings.index_of
      fn (s: string, needle: string) -> i32
      + returns the first byte offset of needle, or -1 when absent
      # strings
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex, string]
      + compiles the pattern for repeated use
      - returns error on invalid pattern
      # regex
    std.regex.match
      fn (re: regex, input: string) -> optional[list[string]]
      + returns captured groups on match
      - returns none otherwise
      # regex
  std.time
    std.time.now_iso_date
      fn () -> string
      + returns today's date in ISO 8601 form
      # time

commit_changelog
  commit_changelog.parse_commit
    fn (raw: string) -> result[commit, string]
    + returns a commit with type, optional scope, subject, and body
    - returns error when the header line cannot be classified
    # parsing
    -> std.strings.split_lines
    -> std.strings.trim
    -> std.regex.match
  commit_changelog.classify
    fn (c: commit) -> string
    + returns a section id such as "features", "fixes", "breaking", or "other"
    + flags a commit as breaking when its body contains a breaking marker
    # classification
    -> std.strings.index_of
    -> std.strings.starts_with
  commit_changelog.group_by_section
    fn (commits: list[commit]) -> map[string, list[commit]]
    + buckets commits into sections in stable order
    # grouping
  commit_changelog.new_template
    fn (sections: list[string]) -> changelog_template
    + creates a template that emits the given sections in order
    # construction
  commit_changelog.render
    fn (tpl: changelog_template, version: string, grouped: map[string, list[commit]]) -> string
    + produces a versioned changelog with a date header and one block per section
    + skips empty sections
    # rendering
    -> std.time.now_iso_date
  commit_changelog.diff_commits
    fn (old_commits: list[commit], new_commits: list[commit]) -> list[commit]
    + returns commits present in new_commits but not in old_commits by subject
    # diff
  commit_changelog.validate
    fn (c: commit) -> result[void, string]
    + accepts commits whose type is one of the recognized set
    - returns error with the offending type name otherwise
    # validation
  commit_changelog.suggest_next_version
    fn (current: string, commits: list[commit]) -> result[string, string]
    + bumps major on any breaking commit, minor on any feature, patch otherwise
    - returns error when current is not a valid semantic version
    # versioning
    -> std.regex.match
