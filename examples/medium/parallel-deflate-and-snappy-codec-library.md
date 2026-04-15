# Requirement: "a parallel deflate and snappy codec library"

Splits input into blocks, encodes or decodes each block through a worker pool, and reassembles in original order.

std
  std.concurrent
    std.concurrent.parallel_map
      fn (items: list[bytes], worker_count: i32, fn: bytes_transform) -> list[result[bytes, string]]
      + applies fn across items using a worker pool, preserving order
      # concurrency
  std.compress
    std.compress.deflate_block
      fn (data: bytes, level: i32) -> bytes
      + compresses a single block with deflate
      # compression
    std.compress.inflate_block
      fn (data: bytes) -> result[bytes, string]
      + decompresses a single deflate block
      - returns error on corrupt input
      # compression
    std.compress.snappy_block_encode
      fn (data: bytes) -> bytes
      + compresses a block with snappy
      # compression
    std.compress.snappy_block_decode
      fn (data: bytes) -> result[bytes, string]
      + decompresses a snappy block
      - returns error on corrupt input
      # compression

pgzp
  pgzp.chunk_input
    fn (data: bytes, block_size: i32) -> list[bytes]
    + splits input into fixed-size blocks
    - returns a single block when data is smaller than block_size
    # chunking
  pgzp.encode_deflate
    fn (data: bytes, block_size: i32, worker_count: i32, level: i32) -> bytes
    + compresses in parallel and concatenates the framed blocks
    # encoding
    -> std.concurrent.parallel_map
    -> std.compress.deflate_block
  pgzp.decode_deflate
    fn (data: bytes, worker_count: i32) -> result[bytes, string]
    + decompresses framed deflate blocks in parallel
    - returns error when a block fails to inflate
    # decoding
    -> std.concurrent.parallel_map
    -> std.compress.inflate_block
  pgzp.encode_snappy
    fn (data: bytes, block_size: i32, worker_count: i32) -> bytes
    + compresses in parallel using the snappy block format
    # encoding
    -> std.concurrent.parallel_map
    -> std.compress.snappy_block_encode
  pgzp.decode_snappy
    fn (data: bytes, worker_count: i32) -> result[bytes, string]
    + decompresses snappy blocks in parallel
    - returns error on corrupt input
    # decoding
    -> std.concurrent.parallel_map
    -> std.compress.snappy_block_decode
