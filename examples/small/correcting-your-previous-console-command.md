# Requirement: "a console command correction suggester"

Given a failed command and its error output, propose a corrected command using a registry of rules.

std: (all units exist)

command_fix
  command_fix.new_registry
    @ () -> registry
    + creates an empty rule registry
    # construction
  command_fix.register_rule
    @ (r: registry, name: string, rule: fix_rule) -> registry
    + adds a named rule; each rule has a matcher and a rewriter
    # extensibility
  command_fix.suggest
    @ (r: registry, command: string, stderr: string, exit_code: i32) -> list[string]
    + returns candidate corrected commands ordered by rule priority
    - returns an empty list when no rule matches
    ? ties are broken by rule registration order
    # suggestion
  command_fix.default_rules
    @ () -> registry
    + returns a registry pre-populated with common corrections such as typo fixes, missing sudo, and wrong subcommand
    # defaults
