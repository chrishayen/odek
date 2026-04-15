# Requirement: "a library for extracting record field names, types, and tags for filtering and exporting"

Given a simple record schema, expose its fields, filter by tag, and produce an export-ready mapping.

std: (all units exist)

fields
  fields.parse_schema
    fn (source: string) -> result[list[field], string]
    + parses a schema definition into a list of fields with name, type, and tag map
    - returns error on malformed syntax
    # parsing
  fields.list_names
    fn (schema: list[field]) -> list[string]
    + returns the field names in declaration order
    # inspection
  fields.filter_by_tag
    fn (schema: list[field], tag_key: string, tag_value: string) -> list[field]
    + returns fields whose tag map has tag_key equal to tag_value
    # filtering
  fields.export_map
    fn (schema: list[field], record: map[string, string]) -> map[string, string]
    + returns a map keyed by each field's export name, using the export tag when present and the field name otherwise
    # export
