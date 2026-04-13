# Requirement: "a tool for generating self-contained mock objects"

Produces a mock type with no runtime dependency on a mocking framework. Each generated method records its invocation and returns caller-supplied stub values.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + returns full file contents as a string
      - returns error when the file cannot be opened
      # filesystem
    std.fs.write_all
      @ (path: string, contents: string) -> result[void, string]
      + writes the contents to the given path, replacing any existing file
      # filesystem
  std.text
    std.text.join
      @ (parts: list[string], sep: string) -> string
      + joins parts with the separator between each pair
      # strings

self_mock
  self_mock.parse_interface
    @ (source: string, interface_name: string) -> result[interface_spec, string]
    + returns a spec describing the interface and each method's parameter and return types
    - returns error when the interface is missing or malformed
    # parsing
  self_mock.render_invocation_struct
    @ (spec: interface_spec) -> string
    + returns source text for a struct that stores one call record per method invocation
    # codegen
    -> std.text.join
  self_mock.render_mock_type
    @ (spec: interface_spec) -> string
    + returns source text for a self-contained mock type with stub setters, call accessors, and no external dependencies
    # codegen
    -> std.text.join
  self_mock.generate_file
    @ (source_path: string, interface_name: string, output_path: string) -> result[void, string]
    + reads, parses, renders, and writes the full self-contained mock file
    - returns error when any step fails
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
