# Requirement: "decode parsed URL query values into typed field values"

Input is an already-parsed key to list-of-values map; the library offers typed accessors.

std: (all units exist)

query_decode
  query_decode.get_string
    @ (values: map[string, list[string]], key: string) -> result[string, string]
    + returns the first value for the key
    - returns error when the key is missing
    # access
  query_decode.get_int
    @ (values: map[string, list[string]], key: string) -> result[i64, string]
    + returns the first value parsed as a signed integer
    - returns error when the key is missing or not an integer
    # access
  query_decode.get_float
    @ (values: map[string, list[string]], key: string) -> result[f64, string]
    + returns the first value parsed as a float
    - returns error when the key is missing or not a float
    # access
  query_decode.get_bool
    @ (values: map[string, list[string]], key: string) -> result[bool, string]
    + returns the first value parsed as a boolean ("1"/"0", "true"/"false")
    - returns error when the key is missing or unrecognized
    # access
  query_decode.get_list_string
    @ (values: map[string, list[string]], key: string) -> list[string]
    + returns all values for the key in order
    + returns an empty list when the key is missing
    # access
