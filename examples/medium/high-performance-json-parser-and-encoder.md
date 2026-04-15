# Requirement: "a high-performance JSON parser and encoder"

Streaming tokenizer feeding a value decoder, plus a value encoder that writes directly to a byte buffer.

std: (all units exist)

jsoniter
  jsoniter.new_tokenizer
    fn (input: bytes) -> tokenizer_state
    + creates a tokenizer positioned at the start of input
    # lexer
  jsoniter.next_token
    fn (state: tokenizer_state) -> result[tuple[json_token, tokenizer_state], string]
    + advances the tokenizer and returns the next token
    + tokens include: begin-object, end-object, begin-array, end-array, colon, comma, string, number, bool, null
    - returns error at unterminated strings or illegal escapes
    - returns error on malformed numeric literals
    # lexer
  jsoniter.parse
    fn (input: bytes) -> result[json_value, string]
    + parses input into a tree of json values
    - returns error when the input contains trailing non-whitespace after the root value
    # parser
    -> jsoniter.next_token
  jsoniter.parse_stream
    fn (input: bytes, visit: fn(path: string, value: json_value) -> void) -> result[void, string]
    + walks input event-style, invoking visit for each leaf value with its JSON pointer
    # parser
    -> jsoniter.next_token
  jsoniter.get_path
    fn (value: json_value, path: string) -> optional[json_value]
    + looks up a value by dot-or-bracket path (e.g. "users[0].name")
    # query
  jsoniter.encode
    fn (value: json_value) -> bytes
    + serializes a json value to a compact byte representation
    + escapes strings according to the JSON specification
    # encoder
  jsoniter.encode_indented
    fn (value: json_value, indent: string) -> bytes
    + serializes with the given indentation string per nesting level
    # encoder
    -> jsoniter.encode
