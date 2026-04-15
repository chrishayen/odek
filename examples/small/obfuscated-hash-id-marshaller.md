# Requirement: "an ID type that marshals to and from an obfuscated hash to avoid exposing raw IDs"

Integer IDs are reversibly obfuscated into short opaque strings using a salt-seeded permutation. The transformation is bijective.

std
  std.encoding
    std.encoding.base36_encode
      fn (value: u64) -> string
      + encodes a 64-bit integer in base36
      # encoding
    std.encoding.base36_decode
      fn (encoded: string) -> result[u64, string]
      + decodes a base36 string to a 64-bit integer
      - returns error on characters outside [0-9a-z]
      # encoding

hideid
  hideid.new_codec
    fn (salt: string, min_length: i32) -> codec
    + creates a codec seeded by the salt; output is padded to at least min_length
    # construction
  hideid.encode
    fn (c: codec, id: u64) -> string
    + returns the obfuscated hash for the id
    + distinct ids produce distinct hashes under the same codec
    # marshal
    -> std.encoding.base36_encode
  hideid.decode
    fn (c: codec, hash: string) -> result[u64, string]
    + returns the original id from a hash produced by the same codec
    - returns error when the hash is malformed
    - returns error when the hash was produced with a different salt
    # unmarshal
    -> std.encoding.base36_decode
