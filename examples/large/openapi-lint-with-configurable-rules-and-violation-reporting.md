# Requirement: "a lightweight linter for OpenAPI documents that runs a configurable set of rules and reports violations"

Load the spec, walk it with a rule registry, collect violations with paths, severities, and messages.

std
  std.json
    std.json.parse
      fn (raw: string) -> result[json_value, string]
      + parses JSON into a dynamic value tree
      - returns error on invalid JSON
      # parsing
  std.yaml
    std.yaml.parse
      fn (raw: string) -> result[json_value, string]
      + parses YAML into the same dynamic value tree as JSON
      - returns error on invalid YAML
      # parsing
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + reads a file into a string
      - returns error when the file cannot be opened
      # filesystem

openapi_lint
  openapi_lint.load
    fn (path: string) -> result[spec_doc, string]
    + reads a spec file (JSON or YAML based on extension)
    - returns error when the file is unreadable or unparseable
    # loading
    -> std.fs.read_all
    -> std.json.parse
    -> std.yaml.parse
  openapi_lint.new_rules
    fn () -> rule_set
    + returns the built-in rule set
    # rules
  openapi_lint.add_rule
    fn (rules: rule_set, id: string, severity: string, check: fn(spec_doc) -> list[violation]) -> rule_set
    + registers an additional rule
    # rules
  openapi_lint.disable_rule
    fn (rules: rule_set, id: string) -> rule_set
    + marks a rule as disabled
    # rules
  openapi_lint.rule_info_required
    fn (spec: spec_doc) -> list[violation]
    + flags missing info.title or info.version
    # rule
  openapi_lint.rule_operation_ids_unique
    fn (spec: spec_doc) -> list[violation]
    + flags duplicate operationId values across paths
    # rule
  openapi_lint.rule_paths_lowercase
    fn (spec: spec_doc) -> list[violation]
    + flags paths containing uppercase letters
    # rule
  openapi_lint.rule_responses_documented
    fn (spec: spec_doc) -> list[violation]
    + flags operations with no documented responses
    # rule
  openapi_lint.run
    fn (rules: rule_set, spec: spec_doc) -> list[violation]
    + runs every enabled rule and concatenates the violations sorted by path
    # execution
  openapi_lint.format_text
    fn (violations: list[violation]) -> string
    + renders violations as human-readable lines with severity and json path
    # reporting
