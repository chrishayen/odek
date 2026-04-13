# Requirement: "a schema-based serialization library for distributed systems"

A schema compiler and a zero-copy binary reader/writer. The project layer parses schemas, generates type info, and encodes/decodes messages.

std: (all units exist)

schema_codec
  schema_codec.parse_schema
    @ (source: string) -> result[schema, string]
    + parses a schema source document into structured type definitions
    - returns error on syntax error
    - returns error when a referenced type is undefined
    # schema
  schema_codec.lookup_struct
    @ (schema: schema, name: string) -> result[struct_def, string]
    + resolves a struct definition by name
    - returns error when the struct is not in the schema
    # schema_lookup
  schema_codec.new_builder
    @ (def: struct_def) -> struct_builder
    + starts a writable message for the given struct
    # encoding
  schema_codec.set_field
    @ (b: struct_builder, field: string, value: field_value) -> result[struct_builder, string]
    + sets a field by name with type checking
    - returns error when the field is undefined on the struct
    - returns error when value type does not match the field type
    # encoding
  schema_codec.freeze
    @ (b: struct_builder) -> bytes
    + serializes the builder into the wire format
    # encoding
  schema_codec.open_reader
    @ (def: struct_def, raw: bytes) -> result[struct_reader, string]
    + creates a reader over raw bytes validated against def
    - returns error when raw is too short for the struct layout
    # decoding
  schema_codec.get_field
    @ (r: struct_reader, field: string) -> result[field_value, string]
    + reads a field by name
    - returns error when the field does not exist
    # decoding
  schema_codec.make_rpc_message
    @ (method: string, payload: bytes) -> bytes
    + frames a payload with method name for transport
    # rpc
