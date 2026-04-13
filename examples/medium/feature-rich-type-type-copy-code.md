# Requirement: "a type-to-type copy code generator"

Given a pair of struct-like type descriptions, emits source code for a function that assigns each field of the source to the matching destination field.

std
  std.fs
    std.fs.read_all
      @ (path: string) -> result[string, string]
      + reads the entire file as a string
      # io
    std.fs.write_all
      @ (path: string, data: string) -> result[void, string]
      + writes the entire file
      # io

convergen
  convergen.parse_type_spec
    @ (source: string) -> result[type_spec, string]
    + parses a struct-like type declaration with named fields and types
    - returns error on unclosed braces or missing field names
    # parsing
  convergen.match_fields
    @ (src: type_spec, dst: type_spec) -> list[field_mapping]
    + returns mappings for every field whose name and type agree
    + case-insensitive name matching
    - fields with incompatible types are skipped
    # mapping
  convergen.convertible
    @ (src_type: string, dst_type: string) -> bool
    + returns true when source type can be assigned or trivially cast to destination
    # mapping
  convergen.emit_function
    @ (name: string, src: type_spec, dst: type_spec, mappings: list[field_mapping]) -> string
    + emits the source of a copy function that assigns each mapped field
    + inserts cast expressions where types differ but are convertible
    # generation
  convergen.generate
    @ (input_path: string, output_path: string, function_name: string) -> result[void, string]
    + reads two type specs from input, writes the generated function to output
    - returns error when parsing either type fails
    # entry
    -> std.fs.read_all
    -> std.fs.write_all
