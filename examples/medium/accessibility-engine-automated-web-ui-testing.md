# Requirement: "an accessibility testing engine that audits a UI tree against a set of rules"

The library takes an abstract node tree (already parsed by the caller) and runs configured checks against it, returning violations. Nothing is framework-specific.

std: (all units exist)

a11y
  a11y.rule_register
    @ (id: string, description: string, check_fn: fn(node) -> list[violation]) -> result[void, string]
    + registers a named rule to be evaluated against nodes
    - returns error when a rule with the same id already exists
    # registry
  a11y.rule_unregister
    @ (id: string) -> result[void, string]
    + removes a previously registered rule
    - returns error when no rule with the id is registered
    # registry
  a11y.audit
    @ (root: node, rule_ids: list[string]) -> result[audit_report, string]
    + walks the tree and runs each named rule against every node
    + returns a report grouping violations by rule id and node path
    - returns error when any requested rule id is not registered
    # audit
    -> a11y.rule_register
  a11y.report_count
    @ (report: audit_report) -> i32
    + returns the total number of violations across all rules
    # report
  a11y.report_by_severity
    @ (report: audit_report, severity: string) -> list[violation]
    + returns all violations of the given severity
    ? severity is one of "minor", "moderate", "serious", "critical"
    # report
