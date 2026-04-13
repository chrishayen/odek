# Requirement: "a resource query language that parses filter, sort, and pagination parameters from an API query string"

Parses a compact filter DSL and common list parameters into a structured query, then translates it into a backend-agnostic query plan.

std
  std.url
    std.url.parse_query
      @ (raw: string) -> map[string, list[string]]
      + parses a url query string into a multimap
      # parsing
  std.json
    std.json.parse_value
      @ (raw: string) -> result[json_value, string]
      + parses a json document into a dynamic value
      - returns error on invalid json
      # serialization

rql
  rql.parse_filter
    @ (expr: string) -> result[filter_expr, string]
    + parses a filter expression like `and(eq(name,"a"),gt(age,21))` into an AST
    - returns error on unbalanced parentheses
    - returns error on unknown operator
    # parsing
    -> std.json.parse_value
  rql.parse_request
    @ (query: string) -> result[query_request, string]
    + parses url query parameters into filter, sort, limit, and offset
    - returns error when limit or offset is not a non-negative integer
    # parsing
    -> std.url.parse_query
  rql.validate
    @ (req: query_request, schema: resource_schema) -> result[void, string]
    + rejects filter, sort, or select references to fields not in the schema
    - returns error listing every unknown field
    # validation
  rql.to_plan
    @ (req: query_request) -> query_plan
    + returns a backend-agnostic plan describing predicates, ordering, and pagination
    # planning
  rql.plan_to_sql
    @ (plan: query_plan, table: string) -> tuple[string, list[string]]
    + renders the plan as a parameterized sql statement and its arguments
    # rendering
