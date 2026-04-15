# Requirement: "a schema emitter that converts struct descriptors into Cap'n Proto schema text"

Accepts a list of struct descriptors and emits a Cap'n Proto schema file body with deterministic field numbering and type mapping.

std: (all units exist)

capnpgen
  capnpgen.new_struct
    fn (name: string) -> capnp_struct
    + creates an empty named struct descriptor
    # construction
  capnpgen.add_field
    fn (s: capnp_struct, name: string, type_name: string) -> capnp_struct
    + appends a field; its slot number is the current field count
    # construction
  capnpgen.map_type
    fn (src_type: string) -> result[string, string]
    + maps a source type name to its Cap'n Proto equivalent (e.g. i32 -> Int32, string -> Text, bytes -> Data, list[T] -> List(T))
    - returns error on unsupported types
    # type_mapping
  capnpgen.generate_id
    fn (name: string) -> u64
    + produces a stable 64-bit file id derived from the top-level name
    # identity
  capnpgen.render_struct
    fn (s: capnp_struct) -> result[string, string]
    + renders a struct block with each field on its own line using its slot number and mapped type
    - returns error if any field has an unmappable type
    # emission
    -> capnpgen.map_type
  capnpgen.render_schema
    fn (file_name: string, structs: list[capnp_struct]) -> result[string, string]
    + renders a full schema file: id header followed by each struct block
    - returns error on the first failing struct
    # emission
    -> capnpgen.generate_id
    -> capnpgen.render_struct
