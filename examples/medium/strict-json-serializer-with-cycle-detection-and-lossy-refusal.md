# Requirement: "a strict JSON serializer that refuses lossy values and detects cycles"

Traverses a generic value tree, rejecting types that do not survive JSON round-trip and returning an error on reference cycles instead of looping forever.

std: (all units exist)

strict_json
  strict_json.encode
    fn (value: json_value) -> result[string, string]
    + serializes null, bool, finite numbers, strings, arrays, and objects
    - returns error when a number is NaN or infinite
    - returns error when a function, date, or regex value is encountered
    - returns error when the value graph contains a cycle
    # entry
  strict_json.walk_detect_cycles
    fn (value: json_value, seen: list[ref_id]) -> result[void, string]
    + traverses the tree and records every container on the path
    - returns error when the same container appears twice on one root-to-leaf path
    # cycle_detection
  strict_json.check_number
    fn (n: f64) -> result[void, string]
    + accepts finite numbers
    - returns error on NaN or infinity
    # numeric_validation
  strict_json.check_kind
    fn (value: json_value) -> result[void, string]
    + accepts only the six JSON-compatible kinds
    - returns error on function, date, regex, or other exotic kinds
    # kind_validation
  strict_json.escape_string
    fn (input: string) -> string
    + returns a JSON-escaped string literal including surrounding quotes
    # escaping
  strict_json.emit_value
    fn (value: json_value) -> result[string, string]
    + produces the JSON text for an already-validated value
    # emit
