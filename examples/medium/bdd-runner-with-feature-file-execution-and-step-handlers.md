# Requirement: "a behavior-driven development test runner that executes feature files against registered step handlers"

Parses Gherkin-style feature text, matches each step line against registered regex handlers, and reports pass or fail per scenario.

std
  std.regex
    std.regex.compile
      @ (pattern: string) -> result[regex, string]
      + compiles the pattern for repeated matching
      - returns error on invalid pattern
      # regex
    std.regex.match
      @ (re: regex, input: string) -> optional[list[string]]
      + returns captured groups when the pattern matches
      - returns none when it does not
      # regex
  std.strings
    std.strings.split_lines
      @ (s: string) -> list[string]
      + splits on newline characters
      # strings
    std.strings.trim
      @ (s: string) -> string
      + removes leading and trailing whitespace
      # strings

bdd_runner
  bdd_runner.parse_feature
    @ (text: string) -> result[feature, string]
    + extracts the feature title and its scenarios, each with ordered step lines
    - returns error when no scenario is present
    # parsing
    -> std.strings.split_lines
    -> std.strings.trim
  bdd_runner.new_suite
    @ () -> suite_state
    + creates an empty suite with no registered handlers
    # construction
  bdd_runner.register_step
    @ (suite: suite_state, pattern: string, handler_id: string) -> result[suite_state, string]
    + adds a handler binding keyed by pattern
    - returns error when the pattern fails to compile
    # registration
    -> std.regex.compile
  bdd_runner.match_step
    @ (suite: suite_state, step_line: string) -> optional[step_match]
    + returns the handler id and captured arguments for the first matching pattern
    - returns none when no pattern matches
    # dispatch
    -> std.regex.match
  bdd_runner.run_scenario
    @ (suite: suite_state, scenario: scenario, invoker: step_invoker) -> scenario_result
    + walks each step, dispatches via the invoker, and records pass or fail
    + stops on the first failing step
    # execution
    -> std.regex.match
  bdd_runner.run_feature
    @ (suite: suite_state, feat: feature, invoker: step_invoker) -> list[scenario_result]
    + runs every scenario in the feature and returns results in order
    # execution
