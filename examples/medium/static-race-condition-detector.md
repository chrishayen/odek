# Requirement: "a library that statically detects race conditions in source code"

Parses a program, builds a shared-state access graph across concurrent tasks, and reports unsynchronized accesses.

std: (all units exist)

chronos
  chronos.parse_program
    fn (source: string) -> result[program_ast, string]
    + parses source code into an AST with function and task boundaries
    - returns error on syntax failure
    # parsing
  chronos.collect_shared_state
    fn (ast: program_ast) -> list[shared_var]
    + identifies variables reachable from more than one concurrent task
    # analysis
  chronos.build_access_graph
    fn (ast: program_ast, vars: list[shared_var]) -> access_graph
    + records each read and write to a shared variable together with its holding locks
    # analysis
  chronos.find_unsynchronized
    fn (graph: access_graph) -> list[race_pair]
    + returns pairs of accesses that can run concurrently without a common lock
    + returns empty list when all accesses are consistently synchronized
    # detection
  chronos.format_report
    fn (pairs: list[race_pair]) -> string
    + formats each race pair with file, line, and the conflicting accesses
    # reporting
