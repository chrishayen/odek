# Requirement: "a FLAC encoder and decoder with support for streaming"

Full round-trip: parse metadata blocks, decode subframes with fixed and LPC predictors, and encode PCM into FLAC frames with a streaming interface.

std
  std.io
    std.io.bit_reader_new
      fn (data: bytes) -> bit_reader_state
      + creates a bit reader positioned at the start of data
      # io
    std.io.bit_reader_read
      fn (state: bit_reader_state, count: i32) -> result[tuple[u64, bit_reader_state], string]
      + reads count bits
      - returns error when fewer than count bits remain
      # io
    std.io.bit_writer_new
      fn () -> bit_writer_state
      + creates an empty bit writer
      # io
    std.io.bit_writer_write
      fn (state: bit_writer_state, value: u64, count: i32) -> bit_writer_state
      + appends the low count bits of value
      # io
    std.io.bit_writer_finish
      fn (state: bit_writer_state) -> bytes
      + pads to a byte boundary and returns the accumulated bytes
      # io
  std.hash
    std.hash.crc8
      fn (data: bytes) -> u8
      + returns the CRC-8 of data using the FLAC polynomial
      # hashing
    std.hash.crc16
      fn (data: bytes) -> u16
      + returns the CRC-16 of data using the FLAC polynomial
      # hashing
    std.hash.md5
      fn (data: bytes) -> bytes
      + returns the 16-byte MD5 digest of data
      # hashing

flac
  flac.decode_stream_header
    fn (data: bytes) -> result[void, string]
    + verifies the four-byte "fLaC" magic at the start of the buffer
    - returns error when the magic is missing
    # metadata
  flac.decode_metadata_blocks
    fn (data: bytes) -> result[tuple[list[metadata_block], i32], string]
    + parses all metadata blocks and returns them with the offset of the first audio frame
    - returns error on truncated blocks or unknown required block types
    # metadata
    -> std.io.bit_reader_read
  flac.decode_frame_header
    fn (reader: bit_reader_state) -> result[tuple[flac_frame_header, bit_reader_state], string]
    + decodes sync, blocking strategy, block size, sample rate, channel assignment, sample size
    - returns error when the header CRC-8 does not match
    # frame
    -> std.io.bit_reader_read
    -> std.hash.crc8
  flac.decode_subframe_fixed
    fn (reader: bit_reader_state, order: i32, block_size: i32, bps: i32) -> result[tuple[list[i32], bit_reader_state], string]
    + decodes a fixed-predictor subframe with the given order
    - returns error on invalid residual partition layout
    # subframe
    -> std.io.bit_reader_read
  flac.decode_subframe_lpc
    fn (reader: bit_reader_state, order: i32, block_size: i32, bps: i32) -> result[tuple[list[i32], bit_reader_state], string]
    + decodes an LPC subframe and applies the prediction filter
    - returns error on invalid QLP precision or coefficient shift
    # subframe
    -> std.io.bit_reader_read
  flac.decode_frame
    fn (reader: bit_reader_state) -> result[tuple[audio_block, bit_reader_state], string]
    + decodes a full frame into interleaved PCM samples
    - returns error when the frame CRC-16 does not match
    # frame
    -> std.hash.crc16
  flac.stream_reader_new
    fn (data: bytes) -> result[stream_reader, string]
    + opens a streaming reader positioned after the metadata blocks
    - returns error when the stream header is invalid
    # streaming
  flac.stream_read_next
    fn (reader: stream_reader) -> result[tuple[audio_block, stream_reader], string]
    + decodes and returns the next frame
    - returns error at end of stream
    # streaming
  flac.encode_frame
    fn (samples: list[i32], sample_rate: i32, channels: i32, bps: i32) -> bytes
    + encodes a single frame using a fixed predictor and appends the CRC-16
    # encoding
    -> std.io.bit_writer_write
    -> std.io.bit_writer_finish
    -> std.hash.crc8
    -> std.hash.crc16
  flac.encode_stream
    fn (blocks: list[audio_block], sample_rate: i32, channels: i32, bps: i32) -> bytes
    + writes the stream header, a STREAMINFO block, and every encoded frame
    ? the final MD5 signature of unencoded samples is stored in STREAMINFO
    # encoding
    -> std.hash.md5
  flac.stream_writer_new
    fn (sample_rate: i32, channels: i32, bps: i32) -> stream_writer
    + opens a streaming writer that buffers samples and flushes frames on demand
    # streaming
  flac.stream_writer_push
    fn (writer: stream_writer, samples: list[i32]) -> tuple[bytes, stream_writer]
    + appends samples and returns any frames that became complete
    # streaming
  flac.stream_writer_finish
    fn (writer: stream_writer) -> bytes
    + flushes the remaining buffered samples as a final frame
    # streaming
