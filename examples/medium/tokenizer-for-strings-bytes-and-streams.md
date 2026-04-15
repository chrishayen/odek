# Requirement: "a tokenizer that can scan strings, byte slices, or a streaming buffer into user-defined tokens"

A lexer with a pluggable token-definition set that can be driven either from a fully-buffered source or a pull-based reader.

std: (all units exist)

tokenizer
  tokenizer.new_definition
    fn () -> token_definition
    + returns an empty token definition set
    # construction
  tokenizer.add_literal
    fn (def: token_definition, kind: string, literal: string) -> token_definition
    + returns a definition that matches the exact literal as the given kind
    # definition
  tokenizer.add_pattern
    fn (def: token_definition, kind: string, match_fn: fn(string, i64) -> i64) -> token_definition
    + returns a definition where match_fn returns the match length at the offset, or -1 when it does not match
    # definition
  tokenizer.scan_string
    fn (def: token_definition, source: string) -> result[list[token], string]
    + returns every token in the source in order
    - returns error when no definition matches at some offset
    # scanning
  tokenizer.scan_bytes
    fn (def: token_definition, source: bytes) -> result[list[token], string]
    + returns every token after decoding source as UTF-8
    - returns error when source is not valid UTF-8
    - returns error when no definition matches at some offset
    # scanning
  tokenizer.new_stream_scanner
    fn (def: token_definition, read_fn: fn() -> optional[string]) -> stream_scanner
    + returns a scanner that pulls more input via read_fn when its buffer is exhausted
    ? read_fn returns absent to signal end of input
    # scanning
  tokenizer.next_token
    fn (scanner: stream_scanner) -> result[optional[token], string]
    + returns the next token, or absent when the stream is exhausted
    - returns error when the current buffer cannot be matched
    # scanning
