# Requirement: "an ordered map type for JSON marshal and unmarshal that preserves key order"

A map-like structure that preserves insertion order, plus JSON encode and decode functions that respect it.

std: (all units exist)

ordered_map_json
  ordered_map_json.new
    fn () -> ordered_map_state
    + returns an empty ordered map
    # construction
  ordered_map_json.set
    fn (state: ordered_map_state, key: string, value: string) -> ordered_map_state
    + inserts or updates a key, appending to the end when new
    # mutation
  ordered_map_json.get
    fn (state: ordered_map_state, key: string) -> optional[string]
    + returns the value when present
    # query
  ordered_map_json.keys
    fn (state: ordered_map_state) -> list[string]
    + returns keys in insertion order
    # iteration
  ordered_map_json.encode
    fn (state: ordered_map_state) -> string
    + encodes as a JSON object with fields emitted in insertion order
    # serialization
  ordered_map_json.decode
    fn (raw: string) -> result[ordered_map_state, string]
    + parses a JSON object preserving the order of keys as they appeared
    - returns error on invalid JSON or non-object root
    # serialization
