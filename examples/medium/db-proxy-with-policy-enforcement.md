# Requirement: "a database protocol-aware proxy that enforces access policies on queries"

Parses each intercepted query enough to classify it, then runs it against a policy engine before forwarding.

std: (all units exist)

db_proxy
  db_proxy.classify_query
    @ (sql: string) -> query_kind
    + returns select/insert/update/delete/ddl based on leading keyword
    - returns unknown for empty or unrecognized statements
    # parsing
  db_proxy.extract_tables
    @ (sql: string) -> list[string]
    + returns table names referenced by the query
    + handles quoted and schema-qualified identifiers
    # parsing
  db_proxy.compile_policy
    @ (source: string) -> result[policy, string]
    + parses a policy document into an evaluable form
    - returns error on malformed rules
    # policy
  db_proxy.evaluate
    @ (policy: policy, principal: string, kind: query_kind, tables: list[string]) -> policy_decision
    + returns allow/deny plus the matching rule name
    + defaults to deny when no rule matches
    # policy
  db_proxy.intercept
    @ (policy: policy, principal: string, sql: string) -> result[string, string]
    + returns the sql unchanged when allowed
    - returns error when the decision is deny
    # enforcement
