# Requirement: "a library to add and list dependencies in a project manifest file"

Reads and edits a TOML-like manifest that lists dependencies under a `[dependencies]` table. The project face is add/list/remove; std handles file IO and manifest parsing.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[string, string]
      + returns the contents of the named file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: string) -> result[void, string]
      + writes the string to the named file, replacing any existing contents
      - returns error when the file cannot be written
      # filesystem
  std.toml
    std.toml.parse
      fn (raw: string) -> result[manifest_ast, string]
      + parses a manifest document into a structured AST preserving tables
      - returns error on malformed syntax
      # parsing
    std.toml.render
      fn (ast: manifest_ast) -> string
      + renders the AST back to a manifest document preserving table order
      # serialization

manifest_deps
  manifest_deps.list
    fn (path: string) -> result[list[tuple[string, string]], string]
    + returns the list of (name, version) pairs under the dependencies table
    - returns error when the manifest does not exist
    # listing
    -> std.fs.read_all
    -> std.toml.parse
  manifest_deps.add
    fn (path: string, name: string, version: string) -> result[void, string]
    + inserts or updates an entry under the dependencies table and writes the file back
    - returns error when the manifest does not contain a dependencies table
    # add_dependency
    -> std.fs.read_all
    -> std.toml.parse
    -> std.toml.render
    -> std.fs.write_all
  manifest_deps.remove
    fn (path: string, name: string) -> result[void, string]
    + removes the named entry and writes the file back
    - returns error when the named dependency is not present
    # remove_dependency
    -> std.fs.read_all
    -> std.toml.parse
    -> std.toml.render
    -> std.fs.write_all
