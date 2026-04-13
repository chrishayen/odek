# Requirement: "a flamegraph generator that resolves sampled program counters using DWARF debug info"

Takes raw stack samples (lists of instruction addresses) and a DWARF-backed symbol table, resolves each address to a function name, aggregates identical stacks, and emits a collapsed-stack representation suitable for a flamegraph renderer.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[bytes, string]
      + reads a file's contents as bytes
      - returns error when the file cannot be opened
      # filesystem

flamegraph
  flamegraph.load_dwarf
    @ (data: bytes) -> result[symbol_table, string]
    + parses DWARF debug sections into a sorted table of address ranges mapped to function names
    - returns error on truncated or malformed DWARF
    # dwarf_parsing
  flamegraph.resolve_address
    @ (st: symbol_table, pc: u64) -> optional[string]
    + returns the function name whose address range contains the program counter
    ? uses binary search over the sorted range table
    # symbol_resolution
  flamegraph.resolve_stack
    @ (st: symbol_table, stack: list[u64]) -> list[string]
    + resolves each program counter in a stack to a function name, filling unknown frames with "[unknown]"
    # symbol_resolution
    -> flamegraph.resolve_address
  flamegraph.new_profile
    @ () -> profile_state
    + creates an empty profile accumulator
    # construction
  flamegraph.record_sample
    @ (p: profile_state, stack: list[string], weight: i64) -> profile_state
    + increments the counter for the given resolved stack by weight
    # aggregation
  flamegraph.fold
    @ (samples: list[tuple[list[u64], i64]], st: symbol_table) -> profile_state
    + resolves and aggregates every raw sample into a profile
    # aggregation
    -> flamegraph.resolve_stack
    -> flamegraph.record_sample
  flamegraph.render_collapsed
    @ (p: profile_state) -> string
    + renders the profile as one line per unique stack: semicolon-joined frames followed by a space and the weight
    # rendering
  flamegraph.generate
    @ (dwarf: bytes, samples: list[tuple[list[u64], i64]]) -> result[string, string]
    + loads DWARF, folds samples, and returns the collapsed-stack text
    - returns error when DWARF parsing fails
    # pipeline
    -> flamegraph.load_dwarf
    -> flamegraph.fold
    -> flamegraph.render_collapsed
