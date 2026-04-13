# Requirement: "a parallel deflate and snappy codec library"

Splits input into blocks, encodes or decodes each block through a worker pool, and reassembles in original order.

std
  std.concurrent
    std.concurrent.parallel_map
      @ (items: list[bytes], worker_count: i32, fn: bytes_transform) -> list[result[bytes, string]]
      + applies fn across items using a worker pool, preserving order
      # concurrency
  std.compress
    std.compress.deflate_block
      @ (data: bytes, level: i32) -> bytes
      + compresses a single block with deflate
      # compression
    std.compress.inflate_block
      @ (data: bytes) -> result[bytes, string]
      + decompresses a single deflate block
      - returns error on corrupt input
      # compression
    std.compress.snappy_block_encode
      @ (data: bytes) -> bytes
      + compresses a block with snappy
      # compression
    std.compress.snappy_block_decode
      @ (data: bytes) -> result[bytes, string]
      + decompresses a snappy block
      - returns error on corrupt input
      # compression

pgzp
  pgzp.chunk_input
    @ (data: bytes, block_size: i32) -> list[bytes]
    + splits input into fixed-size blocks
    - returns a single block when data is smaller than block_size
    # chunking
  pgzp.encode_deflate
    @ (data: bytes, block_size: i32, worker_count: i32, level: i32) -> bytes
    + compresses in parallel and concatenates the framed blocks
    # encoding
    -> std.concurrent.parallel_map
    -> std.compress.deflate_block
  pgzp.decode_deflate
    @ (data: bytes, worker_count: i32) -> result[bytes, string]
    + decompresses framed deflate blocks in parallel
    - returns error when a block fails to inflate
    # decoding
    -> std.concurrent.parallel_map
    -> std.compress.inflate_block
  pgzp.encode_snappy
    @ (data: bytes, block_size: i32, worker_count: i32) -> bytes
    + compresses in parallel using the snappy block format
    # encoding
    -> std.concurrent.parallel_map
    -> std.compress.snappy_block_encode
  pgzp.decode_snappy
    @ (data: bytes, worker_count: i32) -> result[bytes, string]
    + decompresses snappy blocks in parallel
    - returns error on corrupt input
    # decoding
    -> std.concurrent.parallel_map
    -> std.compress.snappy_block_decode
