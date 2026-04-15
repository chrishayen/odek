# Requirement: "an implementation of a compact binary serialization protocol"

Schema-driven binary encoding: a schema declares field names and types, and encode/decode turn records into compact bytes and back.

std: (all units exist)

compactr
  compactr.schema_new
    fn () -> compactr_schema
    + returns an empty schema
    # construction
  compactr.schema_add_field
    fn (s: compactr_schema, name: string, type_tag: string) -> result[compactr_schema, string]
    + returns a new schema with the field appended
    - returns error when type_tag is not one of the known primitive tags
    - returns error when name already exists in the schema
    # schema
  compactr.encode
    fn (s: compactr_schema, record: map[string, bytes]) -> result[bytes, string]
    + returns a compact byte stream: header bitmap for present fields followed by each present field's bytes in schema order
    - returns error when record contains a field not in the schema
    # encoding
  compactr.decode
    fn (s: compactr_schema, data: bytes) -> result[map[string, bytes], string]
    + returns a map of field names to their raw bytes, skipping absent fields
    - returns error when the header bitmap references a field index the schema does not have
    - returns error on truncated input
    # decoding
