# Requirement: "a library to check that version-control commit messages follow configured patterns"

Given a set of named rules (regex patterns) and a range of commits, return per-commit rule matches. Repository reads go through a thin std primitive.

std
  std.vcs
    std.vcs.list_commits
      @ (repo_path: string, from_ref: string, to_ref: string) -> result[list[commit_record], string]
      + returns commit records (hash, author, subject, body) in the inclusive range from..to
      - returns error when a ref does not exist
      # version_control
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[compiled_regex, string]
      + compiles a regex pattern
      - returns error on invalid syntax
      # regex
    std.regex.is_match
      @ (rx: compiled_regex, text: string) -> bool
      + returns true when the pattern matches anywhere in text
      # regex

commit_lint
  commit_lint.load_rules
    @ (rules: list[tuple[string, string]]) -> result[ruleset, string]
    + compiles a list of (name, pattern) entries into a ruleset
    - returns error when any pattern is invalid
    # configuration
    -> std.regex.compile
  commit_lint.check_message
    @ (rules: ruleset, message: string) -> list[string]
    + returns the names of the rules whose patterns do not match the message
    + returns an empty list when all rules match
    # linting
    -> std.regex.is_match
  commit_lint.check_range
    @ (rules: ruleset, repo_path: string, from_ref: string, to_ref: string) -> result[list[tuple[string, list[string]]], string]
    + returns a list of (commit_hash, failing_rule_names) for each commit in the range
    - returns error when the range cannot be read
    # linting
    -> std.vcs.list_commits
