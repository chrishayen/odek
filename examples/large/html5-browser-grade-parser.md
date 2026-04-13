# Requirement: "a browser-grade HTML5 parser"

Tokenizes HTML5 input and builds a DOM tree following the standard insertion-mode state machine. The project owns the state machine; std provides generic character classification.

std
  std.text
    std.text.is_ascii_alpha
      @ (ch: i32) -> bool
      + returns true for A-Z and a-z
      # text
    std.text.ascii_lower
      @ (ch: i32) -> i32
      + maps A-Z to a-z; passes other code points through
      # text
    std.text.decode_numeric_entity
      @ (digits: string, hex: bool) -> result[i32, string]
      + decodes a decimal or hex numeric character reference
      - returns error on out-of-range code points
      # text

html5
  html5.tokenizer_new
    @ (input: string) -> tokenizer_state
    + returns a tokenizer positioned at the start of the input
    # tokenization
  html5.next_token
    @ (state: tokenizer_state) -> tuple[optional[html_token], tokenizer_state]
    + emits start-tag, end-tag, text, comment, and doctype tokens
    + returns (none, final_state) at end of input
    ? the tokenizer is a state machine with data, tag-open, tag-name, and attribute states
    # tokenization
    -> std.text.is_ascii_alpha
    -> std.text.ascii_lower
  html5.decode_named_entity
    @ (name: string) -> optional[string]
    + returns the replacement text for "amp", "lt", "quot", and other named entities
    # entities
  html5.parse_fragment
    @ (input: string, context_element: string) -> result[dom_node, string]
    + parses an HTML fragment inside the given context element
    - returns error when the context element is unknown
    # fragment_parsing
  html5.parse_document
    @ (input: string) -> result[dom_node, string]
    + builds a full document with html, head, and body children
    + recovers from missing or misordered tags per the HTML5 algorithm
    - returns error only on an empty input stream
    # document_parsing
  html5.insertion_mode_step
    @ (parser: parser_state, token: html_token) -> parser_state
    + advances the parser one token using the current insertion mode
    ? insertion modes include initial, before-html, in-head, in-body, after-body
    # insertion_modes
  html5.dom_to_string
    @ (root: dom_node) -> string
    + serializes a DOM tree back to HTML with proper escaping
    # serialization
  html5.get_element_by_id
    @ (root: dom_node, id: string) -> optional[dom_node]
    + walks the tree and returns the first element whose id attribute matches
    # traversal
  html5.get_elements_by_tag
    @ (root: dom_node, tag: string) -> list[dom_node]
    + returns every descendant whose lowercased tag name matches
    # traversal
    -> std.text.ascii_lower
