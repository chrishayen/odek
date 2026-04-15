# Requirement: "a strict code style linter"

A rule-based linter over tokenized source that reports style violations with file location.

std
  std.fs
    std.fs.read_all_text
      fn (path: string) -> result[string, string]
      + returns the file contents as a string
      - returns error when the file does not exist
      # filesystem

style_linter
  style_linter.new_config
    fn () -> lint_config
    + creates a config with all built-in rules enabled at default severity
    # configuration
  style_linter.set_rule
    fn (config: lint_config, rule_name: string, severity: string) -> result[lint_config, string]
    + updates the severity of an existing rule
    - returns error when rule_name is unknown
    - returns error when severity is not one of off, warn, error
    # configuration
  style_linter.tokenize
    fn (source: string) -> list[token]
    + returns the tokens with line and column positions
    ? comments and whitespace are preserved as distinct tokens so whitespace rules can fire
    # tokenization
  style_linter.lint_source
    fn (config: lint_config, path: string, source: string) -> list[violation]
    + runs all enabled rules over the tokens and returns the violations
    + returns an empty list when no rule matches
    # linting
  style_linter.lint_file
    fn (config: lint_config, path: string) -> result[list[violation], string]
    + reads the file and lints its contents
    - returns error when the file cannot be read
    # linting
    -> std.fs.read_all_text
  style_linter.format_violation
    fn (v: violation) -> string
    + returns a "path:line:col severity rule: message" line
    # rendering
