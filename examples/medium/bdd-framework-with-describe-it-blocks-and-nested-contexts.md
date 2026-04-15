# Requirement: "a behavior-driven testing framework with describe/it blocks and nested contexts"

Tests are organized as nested describe contexts containing it cases. The runner walks the tree, executes each case with accumulated setup, and reports results.

std: (all units exist)

bdd
  bdd.describe
    fn (name: string, body: fn(context_builder) -> context_builder) -> context
    + creates a context node with the given name and child nodes from the body
    # structure
  bdd.it
    fn (builder: context_builder, name: string, test: fn() -> test_result) -> context_builder
    + adds a test case to the current context
    # structure
  bdd.before_each
    fn (builder: context_builder, hook: fn() -> void) -> context_builder
    + registers a setup hook run before each test in this context and its descendants
    # hooks
  bdd.after_each
    fn (builder: context_builder, hook: fn() -> void) -> context_builder
    + registers a teardown hook run after each test in this context and its descendants
    # hooks
  bdd.nested
    fn (builder: context_builder, name: string, body: fn(context_builder) -> context_builder) -> context_builder
    + adds a nested describe context
    # structure
  bdd.run
    fn (root: context) -> run_report
    + executes all tests, running accumulated before/after hooks in order
    + returns counts of passed, failed, and the list of failures with dotted names
    # execution
  bdd.format_report
    fn (report: run_report) -> string
    + renders a human-readable summary of the run
    # reporting
