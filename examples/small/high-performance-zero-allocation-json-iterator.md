# Requirement: "a zero-allocation JSON iterator"

Pull-style iterator over a JSON document that yields tokens without materializing an AST.

std: (all units exist)

json_iter
  json_iter.new
    @ (source: string) -> iter_state
    + creates an iterator positioned at the start of source
    # construction
  json_iter.next_token
    @ (it: iter_state) -> result[optional[json_token], string]
    + returns the next structural or value token (object/array start/end, key, string, number, bool, null)
    + returns none at end of input
    - returns error on malformed syntax
    # parsing
  json_iter.skip_value
    @ (it: iter_state) -> result[void, string]
    + advances past the current value, including any nested object or array
    # navigation
    -> json_iter.next_token
  json_iter.read_string_slice
    @ (it: iter_state) -> result[string_span, string]
    + returns start/end offsets into the source for the current string token without copying
    - returns error when the current token is not a string
    # zero_copy
  json_iter.expect
    @ (it: iter_state, kind: token_kind) -> result[void, string]
    + consumes the next token and errors when its kind doesn't match
    # parsing
    -> json_iter.next_token
