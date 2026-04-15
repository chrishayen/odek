# Requirement: "a tool to generate mock implementations of interfaces"

Given an interface definition, emit a mock object whose methods are individually programmable and whose expectations can be verified.

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

mock_gen
  mock_gen.parse_interface
    fn (source: string, interface_name: string) -> result[interface_spec, string]
    + returns a spec with the interface name, method list, and parameter/return types for each method
    - returns error when the interface cannot be located or has malformed signatures
    # parsing
  mock_gen.render_mock
    fn (spec: interface_spec) -> string
    + returns source text for a mock type with per-method expectation setters and a verify helper
    # codegen
  mock_gen.generate_file
    fn (source_path: string, interface_name: string, output_path: string) -> result[void, string]
    + reads the source, parses the interface, renders a mock, and writes it to disk
    - returns error when any step fails
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
