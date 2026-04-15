# Requirement: "a library that analyzes pull-request diffs for security issues and emits review comments"

Ingests a diff, runs a set of pluggable security checks over the changed hunks, and emits review comments targeted at specific lines.

std
  std.strings
    std.strings.split_lines
      fn (s: string) -> list[string]
      + splits on newline characters
      # strings
    std.strings.starts_with
      fn (s: string, prefix: string) -> bool
      + returns true when s begins with prefix
      # strings
    std.strings.contains
      fn (s: string, needle: string) -> bool
      + returns true when needle occurs in s
      # strings
  std.regex
    std.regex.compile
      fn (pattern: string) -> result[regex, string]
      + compiles a pattern for repeated matching
      - returns error on invalid pattern
      # regex
    std.regex.find_all
      fn (re: regex, input: string) -> list[string]
      + returns every match of re in input
      # regex

security_review
  security_review.parse_diff
    fn (raw: string) -> result[list[file_diff], string]
    + returns one file_diff per changed file, each carrying ordered hunks
    - returns error on malformed diff headers
    # parsing
    -> std.strings.split_lines
    -> std.strings.starts_with
  security_review.extract_added_lines
    fn (diff: file_diff) -> list[diff_line]
    + returns only lines prefixed with '+', excluding the file header
    # parsing
  security_review.new_rule_set
    fn () -> rule_set
    + creates an empty rule set
    # construction
  security_review.add_pattern_rule
    fn (set: rule_set, id: string, severity: string, pattern: string, message: string) -> result[rule_set, string]
    + registers a regex-based rule that fires on matching added lines
    - returns error when the pattern does not compile
    # rules
    -> std.regex.compile
  security_review.default_rules
    fn () -> rule_set
    + returns a rule set seeded with common secret and vulnerability patterns
    # rules
    -> std.regex.compile
  security_review.scan_file
    fn (set: rule_set, diff: file_diff) -> list[finding]
    + returns a finding for each rule that matches an added line in the file
    + each finding records the rule id, severity, line number, and snippet
    # analysis
    -> std.regex.find_all
    -> std.strings.contains
  security_review.scan_pull_request
    fn (set: rule_set, diffs: list[file_diff]) -> list[finding]
    + scans every file and concatenates the findings in order
    # analysis
  security_review.finding_to_comment
    fn (f: finding) -> review_comment
    + formats a finding as a review comment with path, line, and message
    # output
  security_review.summarize
    fn (findings: list[finding]) -> map[string, i32]
    + returns counts of findings by severity
    # reporting
  security_review.filter_by_severity
    fn (findings: list[finding], min_severity: string) -> list[finding]
    + returns findings at or above the given severity
    # filtering
  security_review.deduplicate
    fn (findings: list[finding]) -> list[finding]
    + removes duplicate findings by (rule_id, path, line)
    # filtering
