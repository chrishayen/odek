# Requirement: "a library that derives helper functions from type definitions"

Given a declarative type description, generates equality, hashing, and string-conversion helpers as source snippets.

std: (all units exist)

goderive
  goderive.parse_type
    fn (source: string) -> result[type_decl, string]
    + parses a single type declaration into structured fields
    - returns error on malformed input
    # parsing
  goderive.derive_equal
    fn (decl: type_decl) -> string
    + emits source text for a structural equality function over the type's fields
    # derivation
  goderive.derive_hash
    fn (decl: type_decl) -> string
    + emits source text for a deterministic hash over the type's fields
    # derivation
  goderive.derive_string
    fn (decl: type_decl) -> string
    + emits source text for a human-readable string conversion
    # derivation
