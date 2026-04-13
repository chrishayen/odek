# Requirement: "a library for turning source-code identifiers into cross-reference links"

Given a source buffer and a click position, resolve the identifier under the cursor to a canonical symbol URL via a registered resolver.

std: (all units exist)

code_linker
  code_linker.new
    @ () -> linker_state
    + creates a linker with no resolvers registered
    # construction
  code_linker.register_resolver
    @ (state: linker_state, language: string, resolver: symbol_resolver) -> linker_state
    + registers a per-language symbol resolver
    ? symbol_resolver maps (import_path, identifier) to an absolute URL
    # configuration
  code_linker.identifier_at
    @ (source: string, byte_offset: i32) -> result[tuple[string, i32, i32], string]
    + returns the identifier plus start and end byte offsets covering the cursor
    - returns error when the cursor is not over an identifier
    # lexing
  code_linker.resolve_link
    @ (state: linker_state, language: string, source: string, byte_offset: i32) -> result[string, string]
    + returns the URL the identifier under the cursor points to
    - returns error when the cursor is not on an identifier
    - returns error when no resolver is registered for language
    # resolution
