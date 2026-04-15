# Requirement: "a library for producing a filesystem tree diagram from a markdown outline and materializing it as real directories"

Parses an indented markdown list into a node tree, then offers two operations: render the tree as an ASCII diagram, or create the corresponding directories on disk.

std
  std.fs
    std.fs.mkdir_all
      fn (path: string) -> result[void, string]
      + creates the directory and any missing parents
      - returns error when creation fails
      # filesystem
    std.fs.path_join
      fn (base: string, child: string) -> string
      + joins a parent and child path using the host separator
      # filesystem
  std.io
    std.io.write_string
      fn (sink: writer, text: string) -> result[void, string]
      + writes text to a generic writer
      - returns error on write failure
      # io

tree_outline
  tree_outline.parse_markdown
    fn (text: string) -> result[tree_node, string]
    + parses an indented "-"-bullet list into a root tree
    + indent size is inferred from the first indented line
    - returns error when indentation is inconsistent
    # parsing
  tree_outline.render
    fn (root: tree_node) -> string
    + returns the tree as an ASCII diagram using box-drawing characters
    + marks the last child of each parent with the appropriate connector
    # rendering
  tree_outline.write_to
    fn (root: tree_node, sink: writer) -> result[void, string]
    + renders the tree and writes it to the given sink
    - returns error on write failure
    # output
    -> std.io.write_string
  tree_outline.materialize
    fn (root: tree_node, base_path: string) -> result[void, string]
    + creates a directory on disk for every node under base_path
    - returns error when any directory cannot be created
    # filesystem_sync
    -> std.fs.path_join
    -> std.fs.mkdir_all
