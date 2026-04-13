# Requirement: "a typed-interface emitter that converts struct descriptions into a target interface language"

Takes a language-neutral struct description and emits interface declarations in a target type syntax. No source language parsing: the caller supplies already-normalized struct descriptors.

std: (all units exist)

ifacegen
  ifacegen.new_struct
    @ (name: string) -> struct_descriptor
    + creates a named empty struct descriptor
    # construction
  ifacegen.add_field
    @ (s: struct_descriptor, name: string, type_name: string, optional: bool) -> struct_descriptor
    + appends a field with its name, type, and optionality
    # construction
  ifacegen.map_type
    @ (src_type: string) -> string
    + maps a source type name to its target-language equivalent (integer and float types to number, string to string, list[T] to array of mapped T)
    ? unknown types pass through unchanged
    # type_mapping
  ifacegen.emit_interface
    @ (s: struct_descriptor) -> string
    + renders the struct as an interface declaration with each field on its own line
    + marks optional fields with a trailing "?"
    # code_emission
    -> ifacegen.map_type
  ifacegen.emit_all
    @ (structs: list[struct_descriptor]) -> string
    + renders a list of structs as a single concatenated interface block
    # code_emission
    -> ifacegen.emit_interface
