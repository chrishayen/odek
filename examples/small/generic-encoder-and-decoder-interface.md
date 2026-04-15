# Requirement: "a generic interface to encoders and decoders"

A registry that maps a format name to an encode/decode pair so callers can uniformly encode and decode opaque values.

std: (all units exist)

codec
  codec.new_registry
    fn () -> codec_registry
    + returns an empty codec registry
    # construction
  codec.register
    fn (reg: codec_registry, name: string, encode_fn: fn(bytes) -> bytes, decode_fn: fn(bytes) -> result[bytes, string]) -> codec_registry
    + stores the codec pair under the given name, replacing any existing one
    # registration
  codec.encode
    fn (reg: codec_registry, name: string, data: bytes) -> result[bytes, string]
    + returns encoded bytes when the named codec is registered
    - returns error when no codec is registered under the name
    # encoding
  codec.decode
    fn (reg: codec_registry, name: string, data: bytes) -> result[bytes, string]
    + returns decoded bytes when the named codec is registered and decoding succeeds
    - returns error when the codec is missing or decoding fails
    # decoding
