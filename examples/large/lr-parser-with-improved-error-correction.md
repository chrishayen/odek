# Requirement: "an LR parser with improved error correction"

Builds an LR parse table from a grammar, parses token streams, and when parsing fails proposes minimal edits (insertions and deletions) that would make the input accept.

std
  std.collections
    std.collections.priority_queue_new
      fn () -> pqueue_state
      + creates an empty min-priority queue
      # data_structures
    std.collections.priority_queue_push
      fn (q: pqueue_state, item: string, priority: i32) -> pqueue_state
      + inserts an item with a priority key
      # data_structures
    std.collections.priority_queue_pop
      fn (q: pqueue_state) -> tuple[optional[string], pqueue_state]
      + removes and returns the lowest-priority item
      + returns none when the queue is empty
      # data_structures

lr_parser
  lr_parser.parse_grammar
    fn (source: string) -> result[grammar_state, string]
    + parses a BNF-style grammar with rules like "expr: expr '+' term | term"
    - returns error on malformed rules or undefined nonterminals
    # grammar_parsing
  lr_parser.build_item_sets
    fn (g: grammar_state) -> list[item_set]
    + computes the canonical collection of LR(1) item sets
    ? each item set is a closure over kernel items with lookaheads
    # table_construction
  lr_parser.build_table
    fn (g: grammar_state, items: list[item_set]) -> result[parse_table, string]
    + produces action and goto tables from the item sets
    - returns error on shift-reduce or reduce-reduce conflicts
    # table_construction
  lr_parser.tokenize
    fn (g: grammar_state, input: string) -> result[list[token], string]
    + splits the input into tokens according to the grammar's terminal definitions
    - returns error on unrecognized characters
    # lexing
  lr_parser.parse
    fn (table: parse_table, tokens: list[token]) -> result[parse_tree, parse_error]
    + returns a parse tree when the token stream is accepted
    - returns an error with the offending token position and expected set when rejected
    # parsing
  lr_parser.recover_minimal_edits
    fn (table: parse_table, tokens: list[token], err: parse_error) -> list[edit]
    + searches for the smallest sequence of token insertions and deletions that would make the prefix parse
    ? bounded best-first search over edit cost to keep suggestions actionable
    # error_recovery
    -> std.collections.priority_queue_new
    -> std.collections.priority_queue_push
    -> std.collections.priority_queue_pop
  lr_parser.apply_edits
    fn (tokens: list[token], edits: list[edit]) -> list[token]
    + returns a new token stream with the edits applied in order
    # error_recovery
  lr_parser.format_error
    fn (err: parse_error, edits: list[edit]) -> string
    + renders a human-readable diagnostic including the suggested repair
    # diagnostics
