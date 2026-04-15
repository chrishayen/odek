# Requirement: "a performant JSON unmarshaller supporting partial parsing and unknown fields"

Parses JSON into a typed target with the option to tolerate or capture unknown fields and to skip validation for unread sections.

std: (all units exist)

json_unmarshal
  json_unmarshal.parse
    fn (raw: string) -> result[json_value, parse_error]
    + returns a generic value tree for well-formed JSON
    - returns parse_error with byte offset on malformed input
    # parsing
  json_unmarshal.bind
    fn (value: json_value, schema: list[field_spec]) -> result[map[string, json_value], list[string]]
    + returns a map of bound values for each known field in the schema
    - returns the list of field names whose value did not match the declared type
    ? fields not listed in the schema are preserved for later access
    # binding
  json_unmarshal.bind_lax
    fn (value: json_value, schema: list[field_spec]) -> map[string, json_value]
    + returns bound values for every field that converts cleanly and drops the rest
    ? useful when the source has drifted ahead of the consumer
    # binding
  json_unmarshal.unknown_fields
    fn (value: json_value, schema: list[field_spec]) -> map[string, json_value]
    + returns the subset of fields present in value but absent from the schema
    # introspection
