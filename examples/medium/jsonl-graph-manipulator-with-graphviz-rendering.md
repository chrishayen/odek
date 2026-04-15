# Requirement: "a library to manipulate graphs stored as jsonl with graphviz rendering"

Each input line is a JSON object describing a node or an edge. The library loads them into an in-memory graph and can emit DOT source.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      fn (path: string, contents: string) -> result[void, string]
      + writes the contents to the given path, replacing any existing file
      # filesystem
  std.text
    std.text.split_lines
      fn (s: string) -> list[string]
      + splits on newline and drops a trailing empty segment
      # strings
  std.json
    std.json.parse_object
      fn (raw: string) -> result[map[string, string], string]
      + parses a JSON object into a string-to-string map
      - returns error on invalid JSON or non-object root
      # serialization
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

graph_jsonl
  graph_jsonl.load
    fn (raw: string) -> result[graph_state, string]
    + returns a graph built from one JSONL record per line, where records with a "from" field are edges and others are nodes
    - returns error when any line fails to parse as a JSON object
    # loading
    -> std.text.split_lines
    -> std.json.parse_object
  graph_jsonl.add_node
    fn (g: graph_state, id: string, attrs: map[string, string]) -> graph_state
    + returns a graph with the node inserted or its attributes replaced if already present
    # mutation
  graph_jsonl.add_edge
    fn (g: graph_state, from: string, to: string, attrs: map[string, string]) -> graph_state
    + returns a graph with the edge appended
    # mutation
  graph_jsonl.dump_jsonl
    fn (g: graph_state) -> string
    + emits one JSON object per node and edge, terminated by newlines
    # serialization
    -> std.json.encode_object
  graph_jsonl.render_dot
    fn (g: graph_state) -> string
    + returns DOT source that declares every node and edge with their attributes
    # rendering
