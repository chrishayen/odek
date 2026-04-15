# Requirement: "a streaming archive writer supporting zip and tar formats"

Callers push entries one at a time and the library writes archive bytes incrementally to a sink. Format is selected at construction time.

std
  std.compress
    std.compress.deflate
      fn (data: bytes) -> bytes
      + returns the DEFLATE-compressed bytes
      # compression
  std.hash
    std.hash.crc32
      fn (data: bytes) -> u32
      + returns the CRC-32 checksum
      # hash

archiver
  archiver.new
    fn (format: string, sink: byte_sink) -> result[archive_state, string]
    + constructs a "zip" or "tar" writer bound to the sink
    - returns error for unknown formats
    # construction
  archiver.append_entry
    fn (state: archive_state, name: string, content: bytes, mode: u32, mtime: i64) -> result[archive_state, string]
    + writes the framing and compressed payload for the entry
    - returns error when the sink reports failure
    - returns error when name contains ".." path segments
    # append
    -> std.compress.deflate
    -> std.hash.crc32
  archiver.append_stream
    fn (state: archive_state, name: string, chunk_source: byte_source, mode: u32, mtime: i64) -> result[archive_state, string]
    + streams entry data from chunk_source until exhausted
    - returns error on source or sink failure
    # append_stream
    -> std.compress.deflate
    -> std.hash.crc32
  archiver.finish
    fn (state: archive_state) -> result[void, string]
    + writes the central directory or final tar blocks and closes the sink
    - returns error when the sink fails during finalization
    # finalize
