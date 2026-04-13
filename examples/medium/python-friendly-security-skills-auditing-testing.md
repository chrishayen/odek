# Requirement: "a security skills toolkit for auditing, testing, and safer backend development"

A catalog of auditing rules that can be run against code or configuration, with reportable findings at severity levels.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads the full contents of a file
      - returns error when the path does not exist
      # filesystem
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[regex, string]
      + compiles a regex pattern
      - returns error on invalid syntax
      # pattern_matching
    std.regex.find_all
      @ (r: regex, input: string) -> list[string]
      + returns all non-overlapping matches
      # pattern_matching

security_skills
  security_skills.new_catalog
    @ () -> security_catalog
    + creates an empty rule catalog
    # construction
  security_skills.add_pattern_rule
    @ (state: security_catalog, id: string, severity: string, pattern: string, message: string) -> result[security_catalog, string]
    + adds a rule that fires when pattern matches
    - returns error when pattern is not a valid regex
    # rules
    -> std.regex.compile
  security_skills.audit_source
    @ (state: security_catalog, source: string) -> list[finding]
    + runs every rule against the source and returns findings with id, severity, and matched text
    # audit
    -> std.regex.find_all
  security_skills.audit_file
    @ (state: security_catalog, path: string) -> result[list[finding], string]
    + loads the file and audits its contents
    - returns error when the file cannot be read
    # audit
    -> std.fs.read_all
  security_skills.filter_by_severity
    @ (findings: list[finding], min_severity: string) -> list[finding]
    + returns findings at or above the given severity
    ? severities are ordered info < low < medium < high < critical
    # reporting
  security_skills.format_report
    @ (findings: list[finding]) -> string
    + returns a human-readable summary grouped by severity
    # reporting
