# Requirement: "a library that turns source code into a call-graph flowchart in DOT format"

Parses source into function definitions, walks call expressions to build a graph, and emits DOT. The tokenizer is a general-purpose std utility.

std
  std.lex
    std.lex.tokenize_identifiers
      @ (src: string) -> list[token]
      + returns identifier and punctuation tokens with line numbers
      ? skips whitespace, string literals, and comments
      # lexing

code2flow
  code2flow.extract_functions
    @ (src: string, language: string) -> list[function_def]
    + returns one function_def per top-level function or method in src
    - returns empty list when the language is unrecognized
    # parsing
    -> std.lex.tokenize_identifiers
  code2flow.extract_calls
    @ (fn: function_def) -> list[string]
    + returns the names called from within the function body
    # call_extraction
    -> std.lex.tokenize_identifiers
  code2flow.build_graph
    @ (funcs: list[function_def]) -> call_graph
    + creates nodes for every function and edges for every resolved call
    ? unresolved calls are dropped
    # graph_construction
  code2flow.add_module
    @ (graph: call_graph, src: string, language: string) -> call_graph
    + extracts functions and calls from src and merges them into graph
    # graph_merge
  code2flow.to_dot
    @ (graph: call_graph) -> string
    + serializes the graph to GraphViz DOT text
    + escapes special characters in node labels
    # serialization
  code2flow.filter
    @ (graph: call_graph, keep: list[string]) -> call_graph
    + returns a subgraph containing only the named nodes and their edges
    # filtering
