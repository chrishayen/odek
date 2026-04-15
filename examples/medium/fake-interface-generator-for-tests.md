# Requirement: "a tool to generate fake interface implementations for tests"

Parses an interface definition and emits a stub implementation whose methods record calls and return configurable values.

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
    std.text.join
      fn (parts: list[string], sep: string) -> string
      + joins parts with the separator between each pair
      # strings

fake_gen
  fake_gen.parse_interface
    fn (source: string, interface_name: string) -> result[interface_spec, string]
    + returns a spec containing the interface name and an ordered list of methods with their signatures
    - returns error when the named interface is not found in source
    # parsing
  fake_gen.render_fake
    fn (spec: interface_spec) -> string
    + returns source text for a struct that implements every method on the interface
    + each method appends its arguments to a call log and returns configurable stub values
    # codegen
    -> std.text.join
  fake_gen.generate_file
    fn (source_path: string, interface_name: string, output_path: string) -> result[void, string]
    + reads the source file, generates a fake for the named interface, and writes it to the output path
    - returns error when reading, parsing, or writing fails
    # orchestration
    -> std.fs.read_all
    -> std.fs.write_all
