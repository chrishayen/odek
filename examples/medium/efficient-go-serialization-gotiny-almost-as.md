# Requirement: "an efficient binary serialization library"

Encodes typed values into a compact binary format and decodes them back.

std: (all units exist)

tiny
  tiny.new_encoder
    @ () -> encoder_state
    + creates an empty encoder with a growable buffer
    # construction
  tiny.encode_i32
    @ (enc: encoder_state, value: i32) -> encoder_state
    + appends a varint-encoded signed integer
    # encoding
  tiny.encode_i64
    @ (enc: encoder_state, value: i64) -> encoder_state
    + appends a varint-encoded signed integer
    # encoding
  tiny.encode_string
    @ (enc: encoder_state, value: string) -> encoder_state
    + writes length-prefixed UTF-8 bytes
    # encoding
  tiny.encode_bytes
    @ (enc: encoder_state, value: bytes) -> encoder_state
    + writes length-prefixed raw bytes
    # encoding
  tiny.finalize
    @ (enc: encoder_state) -> bytes
    + returns the accumulated encoded bytes
    # encoding
  tiny.new_decoder
    @ (buffer: bytes) -> decoder_state
    + creates a decoder positioned at the start of the buffer
    # construction
  tiny.decode_i32
    @ (dec: decoder_state) -> result[tuple[i32, decoder_state], string]
    + reads a varint-encoded signed integer
    - returns error on truncated input
    # decoding
  tiny.decode_i64
    @ (dec: decoder_state) -> result[tuple[i64, decoder_state], string]
    + reads a varint-encoded signed integer
    - returns error on truncated input
    # decoding
  tiny.decode_string
    @ (dec: decoder_state) -> result[tuple[string, decoder_state], string]
    + reads a length-prefixed UTF-8 string
    - returns error on invalid UTF-8
    # decoding
  tiny.decode_bytes
    @ (dec: decoder_state) -> result[tuple[bytes, decoder_state], string]
    + reads length-prefixed raw bytes
    - returns error on truncated input
    # decoding
