# Requirement: "a library that converts structured records into key-value maps with field selection"

A transformer selects, renames, and omits fields from a source record and emits a map.

std: (all units exist)

structomap
  structomap.new
    @ () -> transformer
    + creates an empty transformer that selects no fields by default
    # construction
  structomap.pick
    @ (t: transformer, field: string) -> transformer
    + adds a field to be copied as-is into the output map
    # selection
  structomap.rename
    @ (t: transformer, source_field: string, target_key: string) -> transformer
    + adds a field to be copied under a different key in the output map
    # renaming
  structomap.omit
    @ (t: transformer, field: string) -> transformer
    + marks a field to be excluded even if previously picked
    # exclusion
  structomap.apply
    @ (t: transformer, record: map[string, string]) -> map[string, string]
    + returns a new map containing the selected, renamed, non-omitted fields
    - returns an empty map when no fields were picked
    # projection
