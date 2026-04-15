# Requirement: "a runtime for data analytics applications"

A typed relational runtime: tabular datasets, columnar storage, a query plan with projection, filter, aggregation, and join, and a small expression language.

std: (all units exist)

analytics
  analytics.new_schema
    fn (columns: list[column_def]) -> schema
    + constructs a schema from typed column definitions
    # schema
  analytics.new_table
    fn (name: string, schema: schema) -> table_state
    + creates an empty table with the given schema
    # construction
  analytics.append_row
    fn (state: table_state, values: list[cell]) -> result[table_state, string]
    + appends a row when values match the schema
    - returns error when count or types disagree with the schema
    # ingestion
  analytics.parse_expr
    fn (source: string) -> result[expr, string]
    + parses a filter or projection expression
    - returns error on syntax error
    # expression_language
  analytics.eval_expr
    fn (e: expr, row: list[cell]) -> result[cell, string]
    + evaluates an expression against a row
    - returns error on type mismatch
    # expression_language
  analytics.plan_scan
    fn (table: string) -> plan_node
    + builds a scan node reading all rows of a table
    # planning
  analytics.plan_project
    fn (input: plan_node, expressions: list[expr]) -> plan_node
    + builds a projection over the input
    # planning
  analytics.plan_filter
    fn (input: plan_node, predicate: expr) -> plan_node
    + builds a filter over the input
    # planning
  analytics.plan_aggregate
    fn (input: plan_node, group_by: list[string], aggs: list[aggregate_spec]) -> plan_node
    + builds a grouped aggregation over the input
    # planning
  analytics.plan_join
    fn (left: plan_node, right: plan_node, left_key: string, right_key: string) -> plan_node
    + builds an equi-join of two inputs
    # planning
  analytics.execute
    fn (plan: plan_node, tables: map[string, table_state]) -> result[table_state, string]
    + executes a plan and returns the result table
    - returns error when a referenced table is missing
    # execution
  analytics.optimize
    fn (plan: plan_node) -> plan_node
    + rewrites the plan, pushing filters below joins and fusing projections
    # optimization
  analytics.count
    fn (table: table_state) -> i64
    + returns the number of rows in a table
    # observability
