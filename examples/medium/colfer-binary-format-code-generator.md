# Requirement: "a code generator for the Colfer binary format"

A library that parses a schema describing typed structs and emits source code for encoders and decoders in a target-language-agnostic intermediate form.

std
  std.fs
    std.fs.read_all
      fn (path: string) -> result[bytes, string]
      + reads the entire contents of a file
      - returns error when the file does not exist
      # filesystem
    std.fs.write_all
      fn (path: string, data: bytes) -> result[void, string]
      + writes bytes to a file, creating or truncating it
      # filesystem

colfer
  colfer.parse_schema
    fn (source: string) -> result[list[struct_decl], string]
    + parses a schema source into a list of struct declarations
    - returns error on unknown field types
    - returns error when field tags are not unique within a struct
    # parsing
  colfer.validate
    fn (decls: list[struct_decl]) -> result[void, string]
    + checks that struct names are unique and referenced types resolve
    - returns error when a field references an unknown struct
    # validation
  colfer.emit_encoder
    fn (decl: struct_decl) -> string
    + returns source for an encoder function that serializes the struct
    ? numeric fields are little-endian, strings are length-prefixed
    # code_generation
  colfer.emit_decoder
    fn (decl: struct_decl) -> string
    + returns source for a decoder function that deserializes the struct
    # code_generation
  colfer.generate
    fn (source: string) -> result[string, string]
    + full pipeline: parses a schema and returns the combined generated source
    - returns error when parsing or validation fails
    # pipeline
  colfer.generate_from_file
    fn (schema_path: string, out_path: string) -> result[void, string]
    + reads a schema file and writes generated source to an output path
    - returns error when the schema file cannot be read
    - returns error when the output file cannot be written
    # pipeline
    -> std.fs.read_all
    -> std.fs.write_all
