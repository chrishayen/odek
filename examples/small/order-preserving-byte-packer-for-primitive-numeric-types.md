# Requirement: "an order-preserving byte packer for primitive numeric types"

Packs signed and unsigned integers into byte sequences whose lexicographic byte-order matches the numeric order of the inputs. Signed integers are biased; unsigned integers use big-endian.

std: (all units exist)

order_pack
  order_pack.pack_u64
    fn (value: u64) -> bytes
    + returns 8 bytes in big-endian order
    + pack_u64(a) < pack_u64(b) lexicographically whenever a < b
    # packing
  order_pack.pack_i64
    fn (value: i64) -> bytes
    + returns 8 bytes, biasing by 2^63 so negative values sort before positive
    + pack_i64(a) < pack_i64(b) lexicographically whenever a < b
    # packing
  order_pack.unpack_u64
    fn (data: bytes) -> result[u64, string]
    + reverses pack_u64
    - returns error when data is not exactly 8 bytes long
    # unpacking
  order_pack.unpack_i64
    fn (data: bytes) -> result[i64, string]
    + reverses pack_i64
    - returns error when data is not exactly 8 bytes long
    # unpacking
