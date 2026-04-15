# Requirement: "a library that builds a json object from key-value pairs"

A single entry point that encodes string key-value pairs as a JSON object.

std
  std.json
    std.json.encode_object
      fn (obj: map[string, string]) -> string
      + encodes a string-to-string map as JSON
      # serialization

json_builder
  json_builder.build
    fn (pairs: list[tuple[string, string]]) -> string
    + returns a JSON object string for the ordered pairs
    ? duplicate keys keep the last value
    # building
    -> std.json.encode_object
