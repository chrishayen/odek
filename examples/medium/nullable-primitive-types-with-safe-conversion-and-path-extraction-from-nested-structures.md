# Requirement: "nullable primitive types with safe conversion and path-based extraction from nested structures"

Typed optional wrappers plus a small path-lookup helper over generic trees.

std: (all units exist)

typ
  typ.none_string
    @ () -> optional[string]
    + returns a null string value
    # construction
  typ.some_string
    @ (s: string) -> optional[string]
    + returns a present string value
    # construction
  typ.to_i64
    @ (v: dyn_value) -> result[i64, string]
    + converts a dynamic value to i64 from int, float, bool, or numeric string
    - returns error when the value is null
    - returns error when a float is not integral
    - returns error when a string is not numeric
    # conversion
  typ.to_f64
    @ (v: dyn_value) -> result[f64, string]
    + converts a dynamic value to f64 from int, float, or numeric string
    - returns error when the value is null
    - returns error when a string is not numeric
    # conversion
  typ.to_string
    @ (v: dyn_value) -> result[string, string]
    + returns the string form of int, float, bool, or string values
    - returns error when the value is null
    # conversion
  typ.to_bool
    @ (v: dyn_value) -> result[bool, string]
    + converts a dynamic value to bool from bool, int (0/1), or "true"/"false"
    - returns error when the value is null
    - returns error on unrecognized string values
    # conversion
  typ.get_path
    @ (root: dyn_value, path: string) -> optional[dyn_value]
    + traverses dotted keys and bracketed indices, e.g. "user.friends[0].name"
    - returns none when any segment is missing
    - returns none when an index is out of bounds
    # extraction
  typ.get_i64
    @ (root: dyn_value, path: string) -> result[i64, string]
    + convenience that combines get_path and to_i64
    - returns error when the path is missing or the value cannot convert
    # extraction
    -> typ.get_path
    -> typ.to_i64
