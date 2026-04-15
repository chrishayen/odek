# Requirement: "a fast struct-to-struct field copier across differing types"

Copies named fields from a source record to a destination record, converting values when the types differ but are compatible.

std: (all units exist)

structcopy
  structcopy.copy_fields
    fn (src: map[string, string], dst: map[string, string]) -> map[string, string]
    + returns a new map with src values overwriting dst entries whose keys match
    + leaves dst entries untouched when the key is not present in src
    # copying
  structcopy.copy_with_rename
    fn (src: map[string, string], dst: map[string, string], rename: map[string, string]) -> map[string, string]
    + copies fields using the rename map to translate src keys to dst keys
    # copying
  structcopy.copy_convert
    fn (src: map[string, string], dst: map[string, string], converters: map[string, string]) -> result[map[string, string], string]
    + applies the named converter ("int", "bool", "trim") to each specified field before copying
    - returns error when a converter fails on a source value
    # conversion
