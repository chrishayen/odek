# Requirement: "a changelog generator driven by a version-control repository with pluggable formatting"

Walks commits between two revisions, groups them by matcher rules, and renders a changelog document.

std
  std.process
    std.process.run
      fn (cmd: string, args: list[string]) -> result[string, string]
      + executes a command and returns stdout
      - returns error on non-zero exit
      # process
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex, string]
      + compiles a regular expression
      - returns error on malformed pattern
      # regex
    std.regex.match_groups
      fn (re: regex, input: string) -> optional[list[string]]
      + returns the captured groups when the input matches
      - returns none when there is no match
      # regex

changelog
  changelog.fetch_commits
    fn (repo_path: string, from_rev: string, to_rev: string) -> result[list[commit], string]
    + returns commits reachable from to_rev but not from_rev, oldest first
    - returns error when either revision is unknown
    # vcs
    -> std.process.run
  changelog.classify
    fn (commits: list[commit], matchers: list[commit_matcher]) -> map[string, list[commit]]
    + groups commits into buckets using the first matching matcher
    + commits that match no matcher are placed in an "other" bucket
    # classification
    -> std.regex.match_groups
  changelog.new_matcher
    fn (label: string, pattern: string) -> result[commit_matcher, string]
    + builds a matcher that buckets commits whose subject matches the pattern
    - returns error on invalid pattern
    # construction
    -> std.regex.compile
  changelog.render_markdown
    fn (groups: map[string, list[commit]], version: string) -> string
    + renders a markdown changelog with one section per group
    + omits empty groups
    # rendering
  changelog.render_plain
    fn (groups: map[string, list[commit]], version: string) -> string
    + renders a plain-text changelog with indented entries per group
    # rendering
