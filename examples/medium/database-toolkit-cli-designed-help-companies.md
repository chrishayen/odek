# Requirement: "a database toolkit for inspecting, diffing, and applying schema changes"

Library-level operations a toolkit would expose: introspect a live schema, diff against a declarative spec, and generate the change plan.

std: (all units exist)

db_toolkit
  db_toolkit.inspect
    @ (executor: sql_executor) -> result[schema, string]
    + returns the current schema by querying information_schema tables
    - returns error when the executor rejects the introspection query
    # introspection
  db_toolkit.parse_hcl_schema
    @ (source: string) -> result[schema, string]
    + parses a declarative schema document into a schema value
    - returns error on unknown block types
    # parsing
  db_toolkit.diff
    @ (current: schema, desired: schema) -> schema_diff
    + returns added/dropped/altered tables and columns
    # diffing
  db_toolkit.plan
    @ (diff: schema_diff) -> list[change]
    + orders changes so dependencies (drops before creates referencing them) come first
    # planning
  db_toolkit.render_sql
    @ (plan: list[change]) -> string
    + renders a plan to an executable SQL script
    # codegen
  db_toolkit.apply
    @ (executor: sql_executor, plan: list[change]) -> result[i32, string]
    + applies each change, returning the number executed
    - returns error on the first failing change
    # application
