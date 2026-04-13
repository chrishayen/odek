# Requirement: "a library that derives helper functions from type definitions"

Given a declarative type description, generates equality, hashing, and string-conversion helpers as source snippets.

std: (all units exist)

goderive
  goderive.parse_type
    @ (source: string) -> result[type_decl, string]
    + parses a single type declaration into structured fields
    - returns error on malformed input
    # parsing
  goderive.derive_equal
    @ (decl: type_decl) -> string
    + emits source text for a structural equality function over the type's fields
    # derivation
  goderive.derive_hash
    @ (decl: type_decl) -> string
    + emits source text for a deterministic hash over the type's fields
    # derivation
  goderive.derive_string
    @ (decl: type_decl) -> string
    + emits source text for a human-readable string conversion
    # derivation
