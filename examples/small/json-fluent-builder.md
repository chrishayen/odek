# Requirement: "a fluent builder for JSON objects"

A chainable builder that accumulates fields and finally serializes. Serialization is a thin std primitive.

std
  std.json
    std.json.encode
      @ (value: json_value) -> string
      + encodes a json value to a JSON string
      # serialization

jsonbuild
  jsonbuild.new_object
    @ () -> json_builder
    + creates an empty object builder
    # construction
  jsonbuild.set
    @ (b: json_builder, key: string, value: json_value) -> json_builder
    + sets a field, replacing any existing value at the key
    # build
  jsonbuild.set_array
    @ (b: json_builder, key: string, values: list[json_value]) -> json_builder
    + sets a field to an array of values
    # build
  jsonbuild.nest
    @ (b: json_builder, key: string, child: json_builder) -> json_builder
    + sets a field to a nested object from another builder
    # build
  jsonbuild.build
    @ (b: json_builder) -> string
    + serializes the builder to a JSON string
    + returns "{}" for an empty builder
    # serialization
    -> std.json.encode
