# Requirement: "a portable dataframe library with a single API across multiple backends"

Build a backend-agnostic expression graph for tables; execute the graph on a chosen backend.

std: (all units exist)

dataframe
  dataframe.table
    fn (name: string, schema: map[string,type_t]) -> expr_node
    + returns a base expression referencing a named table
    # construction
  dataframe.literal
    fn (value: string, ty: type_t) -> expr_node
    + returns a literal scalar expression
    # construction
  dataframe.column
    fn (table: expr_node, name: string) -> result[expr_node, string]
    + returns a column reference from a table expression
    - returns error when the column is not in the schema
    # construction
  dataframe.select
    fn (table: expr_node, columns: list[expr_node]) -> expr_node
    + returns an expression projecting the given columns
    # projection
  dataframe.filter
    fn (table: expr_node, predicate: expr_node) -> expr_node
    + returns an expression that keeps rows where predicate is true
    # filter
  dataframe.join
    fn (left: expr_node, right: expr_node, on: list[tuple[string,string]], kind: join_t) -> expr_node
    + returns a join expression combining left and right on the given column pairs
    # join
  dataframe.group_by
    fn (table: expr_node, keys: list[expr_node], aggregates: list[expr_node]) -> expr_node
    + returns a grouped-aggregation expression
    # aggregation
  dataframe.order_by
    fn (table: expr_node, keys: list[expr_node], ascending: list[bool]) -> expr_node
    + returns an ordered expression
    # order
  dataframe.limit
    fn (table: expr_node, n: i32) -> expr_node
    + returns an expression that keeps the first n rows
    # slice
  dataframe.binary_op
    fn (op: op_t, left: expr_node, right: expr_node) -> expr_node
    + returns a binary expression such as equality, addition, or comparison
    # expression
  dataframe.aggregate
    fn (fn: agg_t, arg: expr_node) -> expr_node
    + returns an aggregate expression such as sum, avg, or count
    # expression
  dataframe.register_backend
    fn (registry: registry_state, name: string, driver: backend_driver) -> registry_state
    + adds a backend driver under the given name
    # registration
  dataframe.new_registry
    fn () -> registry_state
    + returns an empty backend registry
    # construction
  dataframe.compile
    fn (registry: registry_state, backend: string, expr: expr_node) -> result[compiled_plan, string]
    + lowers the expression graph to a backend-specific plan
    - returns error when the backend is unknown
    - returns error when an operation is unsupported by the backend
    # compilation
  dataframe.execute
    fn (registry: registry_state, backend: string, plan: compiled_plan) -> result[result_set, string]
    + runs the compiled plan on the backend and returns rows
    - returns error when the backend driver reports a runtime failure
    # execution
  dataframe.schema_of
    fn (expr: expr_node) -> map[string,type_t]
    + returns the resulting schema of an expression
    # inference
  dataframe.explain
    fn (registry: registry_state, backend: string, expr: expr_node) -> result[string, string]
    + returns a human-readable plan for the expression on the backend
    # debugging
