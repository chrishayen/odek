# Requirement: "coerce values of arbitrary type to a requested target type at runtime"

A dynamic converter: given a source value and a target type descriptor, produce a value of that type using lenient but well-defined rules.

std: (all units exist)

elastic
  elastic.to_i64
    @ (value: dynamic) -> result[i64, string]
    + converts numeric types to i64 with truncation for floats
    + parses numeric strings including hex and binary prefixes
    + converts true to 1 and false to 0
    - returns error for values that have no sensible integer form
    # coercion
  elastic.to_f64
    @ (value: dynamic) -> result[f64, string]
    + converts numeric types to f64
    + parses numeric strings
    - returns error for values that have no sensible float form
    # coercion
  elastic.to_string
    @ (value: dynamic) -> string
    + renders numbers, booleans, and existing strings verbatim
    + renders lists and maps as a compact JSON-like form
    # coercion
  elastic.to_bool
    @ (value: dynamic) -> result[bool, string]
    + treats nonzero numbers, "true"/"yes"/"1" as true
    + treats zero, "false"/"no"/"0" as false
    - returns error for other values
    # coercion
  elastic.to_list
    @ (value: dynamic, element_target: type_desc) -> result[list[dynamic], string]
    + returns a list when the source is already a list, coercing each element
    + wraps a scalar source in a one-element list
    - returns error when element coercion fails
    # coercion
  elastic.to_map
    @ (value: dynamic, key_target: type_desc, val_target: type_desc) -> result[map[dynamic, dynamic], string]
    + returns a map when the source is already a map, coercing keys and values
    - returns error when the source is not a map
    # coercion
  elastic.convert
    @ (value: dynamic, target: type_desc) -> result[dynamic, string]
    + dispatches on the target descriptor to the appropriate specialized coercion
    - returns error when no rule applies
    # dispatch
    -> elastic.to_i64
    -> elastic.to_f64
    -> elastic.to_string
    -> elastic.to_bool
    -> elastic.to_list
    -> elastic.to_map
